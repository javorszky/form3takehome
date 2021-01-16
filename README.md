# Form3 Take Home Exercise

## Instructions
The goal of this exercise is to write a client library in Go to access our fake account API, which is provided as a Docker
container in the file `docker-compose.yaml` of this repository. Please refer to the
[Form3 documentation](http://api-docs.form3.tech/api.html#organisation-accounts) for information on how to interact with the API.

If you encounter any problems running the fake account API we would encourage you to do some debugging first,
before reaching out for help.

### The solution is expected to
- Be written in Go
- Contain documentation of your technical decisions
- Implement the `Create`, `Fetch`, `List` and `Delete` operations on the `accounts` resource. Note that filtering of the List operation is not required, but you should support paging
- Be well tested to the level you would expect in a commercial environment. Make sure your tests are easy to read.

#### Docker-compose
 - Add your solution to the provided docker-compose file
 - We should be able to run `docker-compose up` and see your tests run against the provided account API service 

### Please don't
- Use a code generator to write the client library
- Use (copy or otherwise) code from any third party without attribution to complete the exercise, as this will result in the test being rejected
- Use a library for your client (e.g: go-resty). Only test libraries are allowed.
- Implement an authentication scheme
- Implement support for the fields `data.attributes.private_identification`, `data.attributes.organisation_identification`
  and `data.relationships`, as they are omitted in the provided fake account API implementation
  
## How to submit your exercise
- Include your name in the README. If you are new to Go, please also mention this in the README so that we can consider this when reviewing your exercise
- Create a private [GitHub](https://help.github.com/en/articles/create-a-repo) repository, copy the `docker-compose` from this repository
- [Invite](https://help.github.com/en/articles/inviting-collaborators-to-a-personal-repository) @form3tech-interviewer-1 to your private repo
- Let us know you've completed the exercise using the link provided at the bottom of the email from our recruitment team

## License
Copyright 2019-2021 Form3 Financial Cloud

Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with the License.
You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.

## Implementation by Gabor Javorszky

### Project layout

The entry point is in `cmd/accountclient.go`. In that `main` package I first marshal all the configurations that the application will need, and exit if something is missing / misconfigured. There's no point continuing startup sequence if I know it's not going to work.

Local packages are all withing the `pkg/<pacakgename>` folders.

### Testing

I've been using the https://github.com/stretchr/testify test library for all of my testing and mock generation purposes for the past year. It's served me well, I am comfortable using it, and it makes reading and writing tests more readable as well as allowing me to use convenience assertions like "is this datetime within 15 seconds of this other datetime?" Without the library I would write helper functions to do the checking manually.

### Config package

Separating the config package into its own module allows me to test it in isolation, and gives me the flexibility to add / remove / change what information is passed into the rest of the application, what environment variable keys are used, I can do error checking and validation (make sure a setting that's supposed to be an URL exists, is not empty, is actually an URL).

Normally I would use `spf13/viper` library to offload some of the work needed for that and allow me to parse .env files as well, but due to the limitations of the take home exercise I opted to rewrite a package in simple terms. Because the code is going to run in a docker container, and I can tell docker what environment variables to set, handling `os.GetEnv` and its siblings are enough for this use case.

Note that I have not copy-pasted / adapted `spf13/viper`'s code, merely recreated the same functionality by myself.

### Client package

This is responsible for talking to the test API in the form3 supplied docker image. There's a `New` function that will return a configured Client struct with the base url and the GMT `time.Location` in it. I'm passing in the location because the `New` function should not return an error, which means I had to move functionality that could produce an error outside it. The thinking is that if the application can't create the GMT `time.Location`, it should stop the startup sequence because it won't be able to add the httpdate to the request either way, and there's a bigger problem with the Go runtime in the machine in that case, like failed to download the timezone information, or can't access it on the system.

I've created an `addHeaders` function that would decorate a request, so I don't need to worry about having to add those in each method. This also makes it testable and central, so if I need to fix something, I can do it in one place. Plus it's small, easy to understand.

There's also a helper function that will return the current httpdate in the format needed. Per the [MDN documentation on the Date header](https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Date) the relevant rfc is 7231 section 7.1.1.2, with the format being described in section 7.1.1.1. Go has a builtin time format in the form of `time.RFC1123` which seems to only differ from the one we want in the timezone. The helper function forces the current time to be represented in GMT before being formatted with the RFC1123 format. 

This is also why we need the GMT `time.Location` on the Client struct, so we don't need to create the location each time this helper function is called. 

#### Validation

In the developer documentation for the `Create` endpoint the payloads need to adhere to certain rules based on which country we're trying to add an account to. For this reason I've created client side validation so we don't even send data that would be rejected by the server.

After having written the individual rules and found that most code is repetitive, I've extracted the main validation functionalities into their own functions where I can pass in parameters to check values against.

Two notable exceptions here: Italy's conditional formatting of the bank ID based on whether the account number is present made it necessary to not extract that specific check into a function, as it's not reusable.

The tests cover all documented eventualities.

#### Create

I've made use of `google/uuid` package, because no one should generate uuids by hand. It's possible, but there's no real reason for it. That's needed to generate an ID for a resource that we're creating.

The organisation ID also needs to be an uuidV4, which I will assume comes with the account being used if I were to use an authentication scheme. I have added an example uuidV4 to the environment variables so the config could pick it up and use it throughout the project.
