import gql from "graphql-tag"

export const GET_ACCOUNT = gql`
query getUser($entity: String!) {
  getUser(entity: $entity) {
    id
    entity
    requests
    profile {
      id
      bio
      contact {
        id
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