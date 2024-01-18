import React from 'react'
import { LandingSectionNavContainer, NavbarButton, NavbarContainer, NavbarLogo, NavbarLogoContainer } from './styles'

export default function Navbar() {
  return (
    <NavbarContainer>
        <NavbarLogoContainer>
            <NavbarLogo>KBachStats</NavbarLogo>
        </NavbarLogoContainer>
        <LandingSectionNavContainer>
            <NavbarButton>section1</NavbarButton>
            <NavbarButton>section2</NavbarButton>
            <NavbarButton>section3</NavbarButton>
            <NavbarButton>section4</NavbarButton>
            <NavbarButton>section5</NavbarButton>
        </LandingSectionNavContainer>
        <NavbarButton>Login</NavbarButton>
        <NavbarButton>Register</NavbarButton>
    </NavbarContainer>
  )
}
