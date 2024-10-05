import gql from "graphql-tag"

export const GET_ACCOUNT = gql`
query GetAccount {
  account {
    id
    accountUser {
      id
      username
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