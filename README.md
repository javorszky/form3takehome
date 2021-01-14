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

The entry point is in `cmd/accountclieng.go`. In that `main` package I first marshal all the configurations that the application will need, and exit if something is missing / misconfigured. There's no point continuing startup sequence if I know it's not going to work.

Local packages are all withing the `pkg/<pacakgename>` folders.



### Test package

I've been using the https://github.com/stretchr/testify test library for all of my testing and mock generation purposes for the past year. It's served me well, I am comfortable using it, and it makes reading and writing tests more readable.

### Config package

Separating the config package into its own module allows me to test it in isolation, and gives me the flexibility to add / remove / change what information is passed into the rest of the application, what environment variables keys are used, I can do error checking and validation (make sure a setting that's supposed to be an URL exists, is not empty, is actually an URL).

Normally I would use `spf13/viper` library to offload some of the work needed for that and allow me to parse .env files as well, but due to the limitations of the take home exercise I opted to rewrite a package in simple terms. Because the code is going to run in a docker container, and I can tell docker what environment variables to set, handling `os.GetEnv` and its siblings are enough for this use case.

Note that I have not copy-pasted / adapted `spf13/viper`'s code, merely recreated the same functionality by myself.

###
