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
    width: 360px;
    height: 650px;
    display: flex;
    flex-direction: column;
    align-items: flex-start;
    justify-content: center;
    border: solid 1px white;
    margin:4rem 0;
`

export const RegisterLogo = styled.h1`
    font-size: 1.2rem;
    margin: 0.6rem 2rem;
`

export const RegisterLabel = styled.label`
    display: block;
    margin: 0.3rem 2rem;
`

export const InfoTextDiv = styled.div`
    width: 285px;
    padding: 10px;
    margin: 10px 2rem;
    color: black;
    background: white;
    border-radius: 10px;
    visibility: ${props => props.focus ? "visible" : "hidden"};
`

export const RegisterInput = styled.input`
    margin: 0 2rem;
    padding: 0 0.5rem;
    border: 1px solid white;
    background: none;
    outline: none;
    width: 285px;
    height: 2rem;
    color: rgb(165, 165, 150);
    font-size: 1rem;
    border-radius: 100vh;
`

export const RegisterValidationIcon = styled.span`
    visibility: ${props => props.validInput ? "visible" : "hidden"};
`

export const RegisterInfoText = styled.p`
    font-size: 14px;
`