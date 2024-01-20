import styled from 'styled-components'
import { NavLink } from 'react-router-dom'

export const NavbarContainer = styled.section`
    width: 100%;
    height: 5rem;
    position: fixed;
    border: 1px solid black;
    display: flex;
    flex-direction: row;
    justify-content: right;
    align-items: center;
    z-index: 1;
    background-color: black;
`

export const NavbarButton = styled(NavLink)`
    font-size: 1.1rem;
    margin: 0 2rem 0 0;
    text-decoration: none;
    @media screen and (max-width: 1020px) {font-size: 1rem;}
    color: rgb(195, 165, 110);
    background-color: black;
    cursor: pointer;
`

export const NavbarNavButton = styled.a`
    font-size: 1.1rem;
    margin: 0 2rem 0 0;
    @media screen and (max-width: 1020px) {font-size: 0.9rem;}
    @media screen and (max-width: 860px) {margin: 0 1rem 0 0;}
    color: rgb(165, 165, 150);
    background-color: black;
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
    @media screen and (max-width: 860px) {margin: 1rem;}
    background-color: black;
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
    font-size: 1.4rem;
    padding-left: 2rem;
    color: rgb(165, 165, 150);
`