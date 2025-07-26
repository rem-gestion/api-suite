package main

import (
	"context"
	"fmt"
	"log"
	"time"

	propertypb "github.com/rem-gestion/rem-common/protos/property/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// Connect to the property service
	conn, err := grpc.Dial("localhost:50054", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := propertypb.NewPropertyServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Test ListProperties
	fmt.Println("Testing ListProperties...")
	listResp, err := client.ListProperties(ctx, &propertypb.ListPropertiesRequest{
		Page:  1,
		Limit: 10,
	})
	if err != nil {
		log.Printf("ListProperties failed: %v", err)
	} else {
		fmt.Printf("Found %d properties (total: %d)\n", len(listResp.Properties), listResp.Total)
		for i, property := range listResp.Properties {
			fmt.Printf("  %d. ID: %s, Internal Code: %s\n", i+1, property.Id, property.InternalCode)
		}
	}

	// Test GetProperty (if we have properties)
	if listResp != nil && len(listResp.Properties) > 0 {
		fmt.Printf("\nTesting GetProperty with ID: %s...\n", listResp.Properties[0].Id)
		getResp, err := client.GetProperty(ctx, &propertypb.GetPropertyRequest{
			Id: listResp.Properties[0].Id,
		})
		if err != nil {
			log.Printf("GetProperty failed: %v", err)
		} else {
			fmt.Printf("Property details: ID=%s, Owner=%s, Address=%s\n", 
				getResp.Property.Id, getResp.Property.OwnerPersonId, getResp.Property.AddressId)
		}
	}

	fmt.Println("\ngRPC test completed successfully!")
}
