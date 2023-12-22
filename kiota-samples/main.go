package main

import (
	"os"

	_ "github.com/joho/godotenv/autoload"

	"context"
	"fmt"
	"log"

	"kiota-samples/client"
	"kiota-samples/client/models"

	client2 "kiota-samples/client2"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"

	auth "github.com/microsoft/kiota-abstractions-go/authentication"
	azure "github.com/microsoft/kiota-authentication-azure-go"
	http "github.com/microsoft/kiota-http-go"
)

func main() {
	client_1()
	client_2()
}

func client_1() {
	// API requires no authentication, so use the anonymous
	// authentication provider
	authProvider := auth.AnonymousAuthenticationProvider{}

	// Create request adapter using the net/http-based implementation
	adapter, err := http.NewNetHttpRequestAdapter(&authProvider)
	if err != nil {
		log.Fatalf("Error creating request adapter: %v\n", err)
	}

	// Create the API client
	client := client.NewPostsClient(adapter)

	// GET /posts
	allPosts, err := client.Posts().Get(context.Background(), nil)
	if err != nil {
		log.Fatalf("Error getting posts: %v\n", err)
	}
	fmt.Printf("Retrieved %d posts.\n", len(allPosts))

	// GET /posts/{id}
	specificPostId := int32(5)
	specificPost, err := client.Posts().ByPostIdInteger(specificPostId).Get(context.Background(), nil)
	if err != nil {
		log.Fatalf("Error getting post by ID: %v\n", err)
	}
	fmt.Printf("Retrieved post - ID: %d, Title: %s, Body: %s\n", *specificPost.GetId(), *specificPost.GetTitle(), *specificPost.GetBody())

	// POST /posts
	newPost := models.NewPost()
	userId := int32(42)
	newPost.SetUserId(&userId)
	title := "Testing Kiota-generated API client"
	newPost.SetTitle(&title)
	body := "Hello world!"
	newPost.SetBody(&body)

	createdPost, err := client.Posts().Post(context.Background(), newPost, nil)
	if err != nil {
		log.Fatalf("Error creating post: %v\n", err)
	}
	fmt.Printf("Created new post with ID: %d\n", *createdPost.GetId())

	// PATCH /posts/{id}
	update := models.NewPost()
	newTitle := "Updated title"
	update.SetTitle(&newTitle)

	updatedPost, err := client.Posts().ByPostIdInteger(specificPostId).Patch(context.Background(), update, nil)
	if err != nil {
		log.Fatalf("Error updating post: %v\n", err)
	}
	fmt.Printf("Updated post - ID: %d, Title: %s, Body: %s\n", *updatedPost.GetId(), *updatedPost.GetTitle(), *updatedPost.GetBody())

	// DELETE /posts/{id}
	_, err = client.Posts().ByPostIdInteger(specificPostId).Delete(context.Background(), nil)
	if err != nil {
		log.Fatalf("Error deleting post: %v\n", err)
	}
	fmt.Printf("Deleted post\n")
}

func client_2() {
	clientId := os.Getenv("AZURE_CLIENT_ID")

	// The auth provider will only authorize requests to
	// the allowed hosts, in this case Microsoft Graph
	allowedHosts := []string{"graph.microsoft.com"}
	graphScopes := []string{"User.Read"}

	credential, err := azidentity.NewDeviceCodeCredential(&azidentity.DeviceCodeCredentialOptions{
		ClientID: clientId,
		UserPrompt: func(ctx context.Context, dcm azidentity.DeviceCodeMessage) error {
			fmt.Println(dcm.Message)
			return nil
		},
	})

	if err != nil {
		fmt.Printf("Error creating credential: %v\n", err)
	}

	authProvider, err := azure.NewAzureIdentityAuthenticationProviderWithScopesAndValidHosts(
		credential, graphScopes, allowedHosts)

	if err != nil {
		fmt.Printf("Error creating auth provider: %v\n", err)
	}

	adapter, err := http.NewNetHttpRequestAdapter(authProvider)

	if err != nil {
		fmt.Printf("Error creating request adapter: %v\n", err)
	}

	client := client2.NewGraphApiClient(adapter)

	me, err := client.Me().Get(context.Background(), nil)

	if err != nil {
		fmt.Printf("Error getting user: %v\n", err)
	}

	fmt.Printf("Hello %s, your ID is %s\n", *me.GetDisplayName(), *me.GetId())
}
