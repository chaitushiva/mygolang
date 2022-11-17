package main

import (
	"bytes"
	"crypto/sha1"
	"crypto/tls"
	"flag"
	"fmt"
	"net/url"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/eks"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/shundezhang/oidc-config/pkg/logger"
)

func main() {

	region := flag.String("region", "us-west-2", "AWS Region we need to use")

	profile := flag.String("profile", "x", "AWS profile we need to use")
	CLUSTER := flag.String("CLUSTER", "eks-x", "EKS CLUSTER ID we need to use")

	flag.Parse()
	args := flag.Args()
	fmt.Println(args)
	os.Exit(0)

	fmt.Println(*region)
	fmt.Println(*profile)
	fmt.Println(*CLUSTER)

	sess, err := session.NewSessionWithOptions(session.Options{
		Config: aws.Config{Region: aws.String(*region),
			CredentialsChainVerboseErrors: aws.Bool(true)},
		Profile:           *profile,
		SharedConfigState: session.SharedConfigEnable,
	})
	if err != nil {
		fmt.Println(err)
	}

	svc := eks.New(sess)
	input := &eks.DescribeClusterInput{
		Name: aws.String(*CLUSTER),
	}

	result, err := svc.DescribeCluster(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case eks.ErrCodeResourceNotFoundException:
				fmt.Println(eks.ErrCodeResourceNotFoundException, aerr.Error())
			case eks.ErrCodeClientException:
				fmt.Println(eks.ErrCodeClientException, aerr.Error())
			case eks.ErrCodeServerException:
				fmt.Println(eks.ErrCodeServerException, aerr.Error())
			case eks.ErrCodeServiceUnavailableException:
				fmt.Println(eks.ErrCodeServiceUnavailableException, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return
	}

	OIDC_URL := *result.Cluster.Identity.Oidc.Issuer
	fmt.Println(OIDC_URL)
	CreateOIDCProvider(*profile, *region, OIDC_URL)

}

func CreateOIDCProvider(profile, region, providerUrl string) error {
	log := logger.NewLogger()
	sess, _ := session.NewSessionWithOptions(session.Options{
		Config: aws.Config{Region: aws.String(region),
			CredentialsChainVerboseErrors: aws.Bool(true)},
		Profile:           profile,
		SharedConfigState: session.SharedConfigEnable,
	})
	svc := iam.New(sess)
	input := &iam.CreateOpenIDConnectProviderInput{
		ClientIDList: []*string{
			aws.String("sts.amazonaws.com"),
		},
		ThumbprintList: []*string{
			aws.String(getThumbPrint(providerUrl)),
		},
		Url: aws.String(providerUrl),
	}

	result, err := svc.CreateOpenIDConnectProvider(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case iam.ErrCodeInvalidInputException:
				fmt.Println(iam.ErrCodeInvalidInputException, aerr.Error())
			case iam.ErrCodeEntityAlreadyExistsException:
				fmt.Println(iam.ErrCodeEntityAlreadyExistsException, aerr.Error())
			case iam.ErrCodeLimitExceededException:
				fmt.Println(iam.ErrCodeLimitExceededException, aerr.Error())
			case iam.ErrCodeConcurrentModificationException:
				fmt.Println(iam.ErrCodeConcurrentModificationException, aerr.Error())
			case iam.ErrCodeServiceFailureException:
				fmt.Println(iam.ErrCodeServiceFailureException, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			log.Error(err)
		}
		return err
	}

	fmt.Println(result)
	return nil
}

func getThumbPrint(httpsUrl string) string {
	u, err := url.Parse(httpsUrl)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	add := u.Hostname()
	if u.Port() != "" {
		add = add + ":" + u.Port()
	} else if u.Port() == "" && u.Scheme == "https" {
		add = add + ":443"
	}
	fmt.Printf("getting thumbprint from %s\n", add)
	conn, err := tls.Dial("tcp", add, &tls.Config{})
	if err != nil {
		panic("failed to connect: " + err.Error())
	}

	// Get the ConnectionState struct as that's the one which gives us x509.Certificate struct
	cert := conn.ConnectionState().PeerCertificates[0]
	// fmt.Printf("%s\n", cert.Issuer)
	fingerprint := sha1.Sum(cert.Raw)

	var buf bytes.Buffer
	for _, f := range fingerprint {
		// if i > 0 {
		// 	fmt.Fprintf(&buf, ":")
		// }
		fmt.Fprintf(&buf, "%02X", f)
	}
	fmt.Printf("%x", fingerprint)
	fmt.Printf("Fingerprint for %s: %s", httpsUrl, buf.String())

	defer conn.Close()
	return buf.String()
}
