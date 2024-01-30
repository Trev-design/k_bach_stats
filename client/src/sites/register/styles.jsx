import styled from "styled-components"
import { PageContainer } from "../../styles"

export const RegisterPageContainer = styled(PageContainer)`
    display: flex;
    flex-direction: column;
    width: 100%;
    height: 100%;
    align-items: center;
    justify-content: center;
`

export const RegisterSection = styled.section`
    width: 40vw;
    height: 60vh;
    display: flex;
    flex-direction: column;
    align-items: flex-start;
    justify-content: center;
    border: solid 1px white;
    margin-top: 4rem;
`

export const RegisterLogo = styled.h1`
    font-size: 1.2rem;
    margin: 0.6rem 2rem;
`

export const RegisterLabel = styled.label`
    display: block;
    margin: 0.3rem 2rem;
`

export const RegisterInput = styled.input`
    margin: 0 2rem;
    background: none;
    border-left: none;
    border-right: none;
    border-top: none;
    border-bottom: 1px solid white;
    outline: none;
    width: 30vw;
    height: 2rem;
    color: rgb(165, 165, 150);
    font-size: 1rem;
`

export const RegisterValidationIcon = styled.span`
    visibility: ${props => props.validInput ? "visible" : "hidden"};
`

export const RegisterInfoText = styled.p`
    visibility: ${props => props.valid ? "visible" : "hidden"};
`