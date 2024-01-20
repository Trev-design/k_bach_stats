import React from 'react'
import { LandingSectionNavContainer, NavbarButton, NavbarContainer, NavbarLogo, NavbarLogoContainer, NavbarNavButton } from './styles'

export default function Navbar() {
  return (
    <NavbarContainer>
        <NavbarLogoContainer>
          <NavbarLogo>KBachStats</NavbarLogo>
        </NavbarLogoContainer>
        <LandingSectionNavContainer>
          <NavbarNavButton>section1</NavbarNavButton>
          <NavbarNavButton>section2</NavbarNavButton>
          <NavbarNavButton>section3</NavbarNavButton>
          <NavbarNavButton>section4</NavbarNavButton>
          <NavbarNavButton>section5</NavbarNavButton>
        </LandingSectionNavContainer>
        <NavbarButton to="/login">Login</NavbarButton>
        <NavbarButton to="/register">Register</NavbarButton>
    </NavbarContainer>
  )
}
