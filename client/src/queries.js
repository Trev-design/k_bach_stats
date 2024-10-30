import gql from "graphql-tag"

export const GET_ACCOUNT = gql`
query getUser($userID: String!) {
  getUser(userID: $userID) {
    entity
    requests
    profile {
      bio
      contact {
        name
        email
        imageFilePath
      }
    }
    experiences {
      experience,
      rating,
      id
    }
  }
}
`

export const ADD_EXPERIENCES = gql`
mutation addExperience($input: NewExperienceCredentials!) {
  addExperience(input: $input) 
}
`