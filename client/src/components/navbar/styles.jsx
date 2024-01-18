import styled from 'styled-components'

export const NavbarContainer = styled.section`
    width: 100%;
    height: 5rem;
    position: fixed;
    border: 1px solid black;
    display: flex;
    flex-direction: row;
    justify-content: right;
    align-items: center;
`

export const NavbarButton = styled.a`
    font-size: 1.1rem;
    margin: 0 2rem 0 0;
    cursor: pointer;
`
export const LandingSectionNavContainer = styled.div`
    margin: 0 3rem;
    border: 1px solid black;
    height: 5rem;
    width: 65%;
    display: flex;
    flex-direction: row;
    justify-content: right;
    align-items: center;
`

export const NavbarLogoContainer = styled.div`
    width: 20%;
    height: 5rem;
    display: flex;
    flex-direction: row;
    justify-content: left;
    align-items: center;
`

export const NavbarLogo = styled.h1`
    font-size: 1.2rem;
    padding-left: 2rem;
`