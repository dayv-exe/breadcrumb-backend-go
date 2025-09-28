package helpers

import (
	"breadcrumb-backend-go/models"
	"breadcrumb-backend-go/utils"
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type UserDynamoHelper struct {
	DbClient  *dynamodb.Client
	TableName string
	Ctx       context.Context
}

func (deps *UserDynamoHelper) AddUser(u *models.User, searchTableName string) error {
	// create new nickname item, to reserve nickname in race conditions for example
	newNickname := models.NewNickname(u.Nickname, u.Name, u.Userid)

	// create a slice of add new user db transactions
	newUserTransactions := []types.TransactWriteItem{
		{
			// adds nickname item to reserve name
			Put: &types.Put{
				TableName: aws.String(deps.TableName),
				Item:      *newNickname.DatabaseFormat(),
				// if this fails most likely because nickname is already in use
				ConditionExpression: aws.String("attribute_not_exists(pk)"),
			},
		},
		{
			// add user to db
			Put: &types.Put{
				TableName: aws.String(deps.TableName),
				Item:      *u.DatabaseFormat(),
			},
		},
	}

	searchHelper := SearchDynamoHelper{
		DbClient:        deps.DbClient,
		Ctx:             deps.Ctx,
		SearchTableName: searchTableName,
	}

	// to create all the search indexes for the new user
	// gets all the transactions to add those search indexes to the database
	searchIndexTransactions, siErr := searchHelper.GetUserSearchIndexItems(u)
	if siErr != nil {
		log.Println("error while creating user search indexes")
		return siErr
	}

	input := &dynamodb.TransactWriteItemsInput{
		TransactItems: append(newUserTransactions, searchIndexTransactions...),
	}

	_, err := deps.DbClient.TransactWriteItems(deps.Ctx, input)

	if err != nil {
		// Check for transaction cancellation reasons
		utils.PrintTransactWriteCancellationReason(err)
		return err
	}

	return nil
}

func (deps *UserDynamoHelper) FindByNickname(nickname string) (*models.User, error) {
	input := &dynamodb.QueryInput{
		TableName:              aws.String(deps.TableName),
		IndexName:              aws.String("NicknameIndex"),
		KeyConditionExpression: aws.String("nickname = :nick"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":nick": &types.AttributeValueMemberS{Value: strings.ToLower(nickname)},
		},
		Limit: aws.Int32(1),
	}

	output, err := deps.DbClient.Query(deps.Ctx, input)

	if err != nil {
		return nil, err
	}

	if output.Count == 0 {
		return nil, nil
	}

	return models.ConvertToUser(output.Items[0])
}

func (deps *UserDynamoHelper) FindById(id string) (*models.User, error) {
	input := dynamodb.GetItemInput{
		Key:       models.UserKey(id),
		TableName: &deps.TableName,
	}

	output, err := deps.DbClient.GetItem(deps.Ctx, &input)

	if err != nil {
		return nil, err
	}

	if len(output.Item) == 0 {
		return nil, fmt.Errorf("Unable to find user")
	}

	return models.ConvertToUser(output.Item)
}

func (deps *UserDynamoHelper) DeleteFromDynamo(u *models.User, searchTableName string) error {
	// delete user profile, nickname, friends, post and allat

	// transactions to remove user item and nickname reservation item from db
	deleteUserTransactions := []types.TransactWriteItem{
		{
			// delete account metadata
			Delete: &types.Delete{
				Key:       models.UserKey(u.Userid),
				TableName: aws.String(deps.TableName),
			},
		},
		{
			// delete nickname reservation
			Delete: &types.Delete{
				Key:       models.NicknameKey(u.Nickname),
				TableName: aws.String(deps.TableName),
			},
		},
	}

	searchHelper := SearchDynamoHelper{
		DbClient:        deps.DbClient,
		Ctx:             deps.Ctx,
		SearchTableName: searchTableName,
	}
	// transactions to remove all users search indexes from search table
	searchIndexTransactions, siErr := searchHelper.GetDeleteUserIndexesItems(u)
	if siErr != nil {
		log.Println("error while creating user search indexes")
		return siErr
	}

	input := &dynamodb.TransactWriteItemsInput{
		TransactItems: append(deleteUserTransactions, searchIndexTransactions...),
	}

	_, err := deps.DbClient.TransactWriteItems(deps.Ctx, input)

	if err != nil {
		utils.PrintTransactWriteCancellationReason(err)
		return err
	}

	return nil
}

func (deps *UserDynamoHelper) UpdateName(user *models.User, newName string, updateNickname bool, searchTableName string) error {
	attributeToUpdate := "name"
	// check if the new name is valid
	if !updateNickname {
		// if we are updating name
		if !utils.NameIsValid(&newName) {
			log.Println("New name given to update name in dynamo helper is invalid")
			return fmt.Errorf("New name provided is invalid!")
		}
	}

	if updateNickname {
		// if we are updating nickname
		attributeToUpdate = "nickname"
		nnAvailable, nnaErr := deps.NicknameAvailable(newName)
		if nnaErr != nil {
			// if we are unable to determine if new nickname is available
			log.Println("error while trying to determine if new nickname is available")
			return nnaErr
		}
		// if new nickname is not available
		if !utils.NicknameValid(user.Nickname) || !nnAvailable {
			log.Println("New name given to update nickname in dynamo helper is invalid")
			return fmt.Errorf("New nickname provided is invalid!")
		}
	}

	// will hold all the transactions that need to be carried out to complete the name change
	var transactItems []types.TransactWriteItem

	searchHelper := SearchDynamoHelper{
		DbClient:        deps.DbClient,
		Ctx:             deps.Ctx,
		SearchTableName: searchTableName,
	}

	// get the search indexes of the old name, and puts them all in delete transactions
	indexes, siErr := searchHelper.GetDeleteUserIndexesItems(user)

	if siErr != nil {
		// if an error occurred while trying to get user indexes to delete
		return siErr
	}

	// add them to our transaction slice
	transactItems = append(transactItems, indexes...)

	// create search indexes for the new name, and put them in put transactions
	if !updateNickname {
		user.Name = newName // we can update the user name here after we delete the indexes with the old name
	} else {
		user.Nickname = newName // we can update the user nickname here after we delete the indexes with the old nickname
	}
	indexes, siErr = searchHelper.GetUserSearchIndexItems(user)
	if siErr != nil {
		// if an error occurred while trying to get user indexes to delete
		return siErr
	}

	// add the new name transactions to our transaction slice
	transactItems = append(transactItems, indexes...)

	// the transaction to actually update the name
	userItems := types.TransactWriteItem{
		Update: &types.Update{
			Key:              models.UserKey(user.Userid),
			TableName:        aws.String(deps.TableName),
			UpdateExpression: aws.String(fmt.Sprintf("SET %s = :n", attributeToUpdate)),
			ExpressionAttributeValues: map[string]types.AttributeValue{
				":n": &types.AttributeValueMemberS{Value: newName},
			},
		},
	}

	// add the name update transaction to the transaction slice
	transactItems = append(transactItems, userItems)

	input := &dynamodb.TransactWriteItemsInput{
		TransactItems: transactItems,
	}

	_, updateErr := deps.DbClient.TransactWriteItems(deps.Ctx, input)
	if updateErr != nil {
		log.Println("Something went wrong while trying to update name")
		utils.PrintTransactWriteCancellationReason(updateErr)
		return updateErr
	}

	return nil
}

func (deps *UserDynamoHelper) updateAttribute(attributeName string, key map[string]types.AttributeValue, newAttribute types.AttributeValue) error {
	input := &dynamodb.UpdateItemInput{
		Key:                 key,
		TableName:           &deps.TableName,
		ConditionExpression: aws.String(fmt.Sprintf("SET %s = :s", attributeName)),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":s": newAttribute,
		},
	}

	_, err := deps.DbClient.UpdateItem(deps.Ctx, input)
	if err != nil {
		log.Print("An error occurred while trying to update attribute")
		return err
	}

	return nil
}

func (deps *UserDynamoHelper) UpdateBio(userId string, bio string) error {
	if !utils.BioIsValid(&bio) {
		log.Println("new bio given to update bio function is invalid")
		return fmt.Errorf("Bio invalid!")
	}

	key := models.UserKey(userId)
	val := &types.AttributeValueMemberS{Value: bio}

	return deps.updateAttribute("bio", key, val)
}

func (deps *UserDynamoHelper) UpdateDpUrl(userId string, url string) error {
	return nil
}

func (deps *UserDynamoHelper) NicknameAvailable(nickname string) (bool, error) {
	input := dynamodb.GetItemInput{
		Key:       models.NicknameKey(nickname),
		TableName: aws.String(deps.TableName),
	}

	return nicknameAvailableQueryRunner(func() (*dynamodb.GetItemOutput, error) {
		return deps.DbClient.GetItem(deps.Ctx, &input)
	})
}

func nicknameAvailableQueryRunner(queryFn func() (*dynamodb.GetItemOutput, error)) (bool, error) {
	result, err := queryFn()

	if err != nil {
		return false, err
	}

	isAvailable := result.Item == nil
	return isAvailable, nil
}
