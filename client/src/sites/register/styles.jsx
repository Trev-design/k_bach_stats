import styled from "styled-components"
import { PageContainer } from "../../styles"
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome'

export const RegisterPageContainer = styled(PageContainer)`
    display: flex;
    flex-direction: column;
    width: 100vw;
    height: 100vh;
    align-items: center;
    justify-content: center;
    border: 1px solid white;
`

export const RegisterLabel = styled.label`
    display: block;
`

export const RegisterValidationIcon = styled.span`
    visibility: ${props => props.validInput ? "visible" : "hidden"};
`

export const RegisterInfoText = styled.p`
    visibility: ${props => props.valid ? "visible" : "hidden"};
`