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
    workspaces {
      id
      name
    }
  }
}
`