package main

import (
	"fmt"
	"os"

	awscdk "github.com/aws/aws-cdk-go/awscdk/v2"
	apigateway "github.com/aws/aws-cdk-go/awscdk/v2/awsapigateway"
	awslambda "github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	constructs "github.com/aws/constructs-go/constructs/v10"
	jsii "github.com/aws/jsii-runtime-go"
)

type GoLambdaCdkStackProps struct {
	awscdk.StackProps
}

func NewGoLambdaCdkStack(scope constructs.Construct, id string, props *GoLambdaCdkStackProps) awscdk.Stack {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}
	stack := awscdk.NewStack(scope, &id, &sprops)

	// Lambda function
	lambdaPath, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	lambdaPath = fmt.Sprintf("%s/lambdas", lambdaPath)

	lambdaFunction := awslambda.NewFunction(stack, jsii.String("GoLambdaFunction"), &awslambda.FunctionProps{
		// Runtime: awslambda.Runtime_GO_1_X(),
		Runtime: awslambda.Runtime_PROVIDED_AL2023(),
		Handler: jsii.String("main"),
		Code:    awslambda.Code_FromAsset(jsii.String(lambdaPath), nil),
		Environment: &map[string]*string{
			"ENV": jsii.String("DEV"),
		},
	})

	// API Gateway
	api := apigateway.NewRestApi(stack, jsii.String("GoLambdaApi"), &apigateway.RestApiProps{
		RestApiName: jsii.String("GoLambdaApi"),
		Description: jsii.String("API Gateway for Go Lambda"),
		// Enable CORS
		DefaultCorsPreflightOptions: &apigateway.CorsOptions{
			AllowOrigins: jsii.Strings("*"),                      // Add your allowed origins
			AllowMethods: jsii.Strings("GET", "POST", "OPTIONS"), // Only the methods you need
			AllowHeaders: jsii.Strings("Content-Type"),           // Only the headers you need, if there are other headers in the lambda, add them here
		},
	})

	// Integrate Lambda with API Gateway
	lambdaIntegration := apigateway.NewLambdaIntegration(lambdaFunction, &apigateway.LambdaIntegrationOptions{})

	// Create a resource (e.g., /hello)
	resource := api.Root().AddResource(jsii.String("hello"), &apigateway.ResourceOptions{})
	// Add a GET method to the resource that points to the Lambda function
	resource.AddMethod(jsii.String("GET"), lambdaIntegration, &apigateway.MethodOptions{})
	// create a post method
	resource.AddMethod(jsii.String("POST"), lambdaIntegration, &apigateway.MethodOptions{})

	return stack
}

func main() {
	defer jsii.Close()

	app := awscdk.NewApp(nil)

	NewGoLambdaCdkStack(app, "GoLambdaCdkStack", &GoLambdaCdkStackProps{
		awscdk.StackProps{
			Env: env(),
		},
	})

	app.Synth(nil)
}

// env determines the AWS environment (account+region) in which our stack is to
// be deployed. For more information see: https://docs.aws.amazon.com/cdk/latest/guide/environments.html
func env() *awscdk.Environment {
	// If unspecified, this stack will be "environment-agnostic".
	// Account/Region-dependent features and context lookups will not work, but a
	// single synthesized template can be deployed anywhere.
	//---------------------------------------------------------------------------
	// return nil

	// Uncomment if you know exactly what account and region you want to deploy
	// the stack to. This is the recommendation for production stacks.
	//---------------------------------------------------------------------------
	// return &awscdk.Environment{
	//  Account: jsii.String("123456789012"),
	//  Region:  jsii.String("us-east-1"),
	// }

	// Uncomment to specialize this stack for the AWS Account and Region that are
	// implied by the current CLI configuration. This is recommended for dev
	// stacks.
	//---------------------------------------------------------------------------
	return &awscdk.Environment{
		Account: jsii.String(os.Getenv("CDK_DEFAULT_ACCOUNT")),
		Region:  jsii.String(os.Getenv("CDK_DEFAULT_REGION")),
	}
}
