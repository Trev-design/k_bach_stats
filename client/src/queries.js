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

export const ADD_EXPERIENCE = gql`
mutation addExperience($input: ExperienceCredentials!) {
  addExperience(input: $input) {
    experience
    rating
    id
  }
}
`

export const ADD_NEW_EXPERIENCE = gql`
mutation createExperience($input: NewExperienceCreddentials!) {
  createExperience(input: $input) {
    experience
    rating
    id
  }
}
`

export const ADD_EXPERINECE_BATCH = gql`
mutation batchNewExperiences($existing: [ExperienceCredentials!]!, $new: [NewExperienceCredentials!]!) {
  batchNewExperiences(existing: $existing, new: $new) {
    experience
    rating
    id
  }
}
`