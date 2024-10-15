import gql from "graphql-tag"

export const GET_ACCOUNT = gql`
query GetAccount($entity: String!) {
  account(entity: $entity) {
    id
    accountUser {
      id
      profile {
        description
        contact {
          name
          email
          profileImagepath
        }
      }
    }
    workSpaces {
      name
      id
    }
  }
}
`