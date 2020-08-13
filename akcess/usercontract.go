package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// UserContract contract for storing user in blockchain
type UserContract struct {
	contractapi.Contract
}

// CreateUser adds a new user to the world state with given details
func (u *UserContract) CreateUser(ctx contractapi.TransactionContextInterface) (string, error) {
	akcessid, _ := ctx.GetClientIdentity().GetID()
	userAsBytes, err := ctx.GetStub().GetState(akcessid)
	txID := ctx.GetStub().GetTxID()

	if err != nil {
		return txID, fmt.Errorf("Failed to read from world state. %s", err.Error())
	}
	if userAsBytes != nil {
		return txID, fmt.Errorf("AKcessID %s already exist", akcessid)
	}

	user := User{
		ObjectType:    "user",
		AkcessID:      akcessid,
		Verifications: map[string][]Verification{},
	}

	newUserAsBytes, _ := json.Marshal(user)

	fmt.Printf("%s: User with AKcessID %s added\n", txID, akcessid)
	return txID, ctx.GetStub().PutState(akcessid, newUserAsBytes)
}

// CreateVerifier register new verifier in Blockchain
func (u *UserContract) CreateVerifier(ctx contractapi.TransactionContextInterface, verifierName string, VerifierGrade string) (string, error) {
	akcessid, _ := ctx.GetClientIdentity().GetID()
	verifierAsBytes, err := ctx.GetStub().GetState(akcessid)
	txID := ctx.GetStub().GetTxID()

	if err != nil {
		return txID, fmt.Errorf("Failed to read from world state. %s", err.Error())
	}
	if verifierAsBytes != nil {
		return txID, fmt.Errorf("AKcessID %s already exist", akcessid)
	}

	verifier := Verifier{
		ObjectType:    "verifier",
		AkcessID:      akcessid,
		VerifierName:  verifierName,
		VerifierGrade: VerifierGrade,
	}

	newVerifierAsBytes, _ := json.Marshal(verifier)

	fmt.Printf("%s: Verifier with AKcessID %s added\n", txID, akcessid)
	return txID, ctx.GetStub().PutState(akcessid, newVerifierAsBytes)
}

// AddUserProfileVerification add verifcation transaction and field of users profiles is verfiied
func (u *UserContract) AddUserProfileVerification(ctx contractapi.TransactionContextInterface, userAKcessID string, profileField string, expiryDate string) (string, error) {
	verifierAKcessID, _ := ctx.GetClientIdentity().GetID()
	txID := ctx.GetStub().GetTxID()
	verifierAsBytes, err := ctx.GetStub().GetState(verifierAKcessID)
	if err != nil {
		return txID, fmt.Errorf("Failed to read from world state. %s", err.Error())
	}
	if verifierAsBytes == nil {
		return txID, fmt.Errorf("AKcessID %s doesn't exist", verifierAKcessID)
	}

	// if !IsVerifier(ctx) {
	// 	return txID, fmt.Errorf("Person who is invoking a transaction is not a verifier")
	// }

	expirydate, err := time.Parse(time.RFC3339, expiryDate)
	if err != nil {
		panic(err)
	}

	var verifier Verifier
	json.Unmarshal(verifierAsBytes, &verifier)

	verification := Verification{
		VerifierObj: verifier,
		ExpirtyDate: expirydate,
	}

	userAsBytes, err := ctx.GetStub().GetState(userAKcessID)
	if err != nil {
		return txID, fmt.Errorf("Failed to read from world state. %s", err.Error())
	}

	if userAsBytes == nil {
		return txID, fmt.Errorf("AKcessID %s doesn't exist", userAKcessID)
	}

	var user User
	json.Unmarshal(userAsBytes, &user)

	verifierList := VerifiersList(user.Verifications[profileField])

	_, found := Find(verifierList, verifierAKcessID)
	if found {
		for i, v := range user.Verifications[profileField] {
			if v.VerifierObj.AkcessID == verifierAKcessID {
				user.Verifications[profileField][i].ExpirtyDate = expirydate
				break
			}
		}
	} else {
		user.Verifications[profileField] = append(user.Verifications[profileField], verification)
	}

	userAsBytes, _ = json.Marshal(user)

	fmt.Printf("%s: Profile field %s of user %s verified by %s\n", txID, profileField, userAKcessID, verifierAKcessID)
	return txID, ctx.GetStub().PutState(userAKcessID, userAsBytes)
}

// GetVerifiersOfUserProfile get verifiers of perticular user field
func (u *UserContract) GetVerifiersOfUserProfile(ctx contractapi.TransactionContextInterface, akcessid string, profileField string) ([]Verification, error) {
	userAsBytes, err := ctx.GetStub().GetState(akcessid)

	if err != nil {
		return nil, fmt.Errorf("Failed to read from world state. %s", err.Error())
	}
	if userAsBytes == nil {
		return nil, fmt.Errorf("AKcessID %s doesn't exist", akcessid)
	}

	var user User
	json.Unmarshal(userAsBytes, &user)
	// verificationList := VerifiersList(user.Verifications[profileField])

	return user.Verifications[profileField], nil
}

// GetVerifier get verifier
func (u *UserContract) GetVerifier(ctx contractapi.TransactionContextInterface, akcessid string) (*Verifier, error) {
	verifierAsBytes, err := ctx.GetStub().GetState(akcessid)

	if err != nil {
		return nil, fmt.Errorf("Failed to read from world state. %s", err.Error())
	}
	if verifierAsBytes == nil {
		return nil, fmt.Errorf("AKcessID %s doesn't exist", akcessid)
	}

	var verifier Verifier
	json.Unmarshal(verifierAsBytes, &verifier)

	return &verifier, nil
}
