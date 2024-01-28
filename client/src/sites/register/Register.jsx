import React from 'react'
import {RegisterInfoText, RegisterLabel, RegisterPageContainer, RegisterValidationIcon } from './styles'
import { useRef, useState, useEffect } from 'react'
import {faInfoCircle, faCheck, faTimes} from '@fortawesome/free-solid-svg-icons'
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome'

function Register() {
  const EMAIL_REGEX = /^[A-Za-z0-9._%+-§$%!?]+@[A-Za-z0-9.-]+\.[a-z]{2,4}$/
  const PASSWORD_REGEX = /^(?=.*[a-z])(?=.*[A-Z])(?=.*[0-9])(?=.*[!?&%$§#@€]).{10,}$/

  const userRef = useRef()
  const errRef = useRef()

  const [name, setName] = useState('')
  const [validName, setValidName] = useState(false)
  const [nameFocus, setNamefocus] = useState(false)

  const [user, setUser] = useState('')
  const [validEmail, setValidEmail] = useState(false)
  const [userFocus, setUserfocus] = useState(false)

  const [password, setPassword] = useState('')
  const [validPassword, setValidPassword] = useState(false)
  const [passworFocus, setPasswordFocus] = useState(false)

  const [confirmation, setConfirmation] = useState('')
  const [validConfirmation, setValidConfirmation] = useState(false)
  const [confirmationFocus, setConfirmationFocus] = useState(false)

  const [errorMessage, setErrorMessage] = useState('')
  const [success, setSuccess] = useState(false)

  useEffect(() => {
    userRef.current.focus()
  }, [])

  useEffect(() => {
    setValidName(name.length >= 4)
  })

  useEffect(() => {
    const result = EMAIL_REGEX.test(user)
    setValidEmail(result)
  }, [user])

  useEffect(() => {
    const result = PASSWORD_REGEX.test(password)
    setValidPassword(result)
    setValidConfirmation(password === confirmation)
  }, [password, confirmation])

  useEffect(() => {
    setErrorMessage('')
  }, [user, password, confirmation])

  return (
    <RegisterPageContainer>
      <section>
      <p ref={errRef} aria-aria-live='assertlive'>{errorMessage}</p>
      <h1> Register </h1>
      <form action="">

      <RegisterLabel htmlFor="email">
          Email:
          <RegisterValidationIcon validInput={validName && name !== ''}>
            <FontAwesomeIcon icon={faCheck} color='#38FF5D'/>
          </RegisterValidationIcon>
          <RegisterValidationIcon validInput={!validName && name !== ''}>
            <FontAwesomeIcon icon={faTimes} color='#FF6161'/>
          </RegisterValidationIcon>
        </RegisterLabel>
        <input 
          type="text"
          id='email'
          ref={userRef}
          autoComplete='off'
          onChange={e => setName(e.target.value)}
          required
          aria-invalid={validName ? "false" : "true"}
          aria-describedby='uidnote'
          onFocus={() => setNamefocus(true)}
          onBlur={() => setNamefocus(false)}
        />
        <RegisterInfoText id='uidnote' valid={!validName && name !== '' && nameFocus}>
          <FontAwesomeIcon icon={faInfoCircle}/>
          invalid email
        </RegisterInfoText>

        <RegisterLabel htmlFor="email">
          Email:
          <RegisterValidationIcon validInput={validEmail && user !== ''}>
            <FontAwesomeIcon icon={faCheck} color='#38FF5D'/>
          </RegisterValidationIcon>
          <RegisterValidationIcon validInput={!validEmail && user !== ''}>
            <FontAwesomeIcon icon={faTimes} color='#FF6161'/>
          </RegisterValidationIcon>
        </RegisterLabel>
        <input 
          type="text"
          id='email'
          ref={userRef}
          autoComplete='off'
          onChange={e => setUser(e.target.value)}
          required
          aria-invalid={validEmail ? "false" : "true"}
          aria-describedby='uidnote'
          onFocus={() => setUserfocus(true)}
          onBlur={() => setUserfocus(false)}
        />
        <RegisterInfoText id='uidnote' valid={!validEmail && user !== '' && userFocus}>
          <FontAwesomeIcon icon={faInfoCircle}/>
          invalid email
        </RegisterInfoText>

        <RegisterLabel htmlFor="password">
          Password:
          <RegisterValidationIcon validInput={validPassword && password !== ''}>
            <FontAwesomeIcon icon={faCheck} color='#38FF5D'/>
          </RegisterValidationIcon>
          <RegisterValidationIcon validInput={!validPassword && password !== ''}>
            <FontAwesomeIcon icon={faTimes} color='#FF6161'/>
          </RegisterValidationIcon>
        </RegisterLabel>
        <input 
          type="password"
          id='password'
          ref={userRef}
          autoComplete='off'
          onChange={e => setPassword(e.target.value)}
          required
          aria-invalid={validPassword ? "false" : "true"}
          aria-describedby='pwdnote'
          onFocus={() => setPasswordFocus(true)}
          onBlur={() => setPasswordFocus(false)}
        />
        <RegisterInfoText id='pwdnote'valid={!validPassword && password !== '' && passworFocus}>
          <FontAwesomeIcon icon={faInfoCircle}/>
          password must have at least 10 characters and lowercase and uppercase letters. also it needs to have digits and special characters like %&$§!?#€&
        </RegisterInfoText>

        <RegisterLabel htmlFor="confirm">
          Confirmation:
          <RegisterValidationIcon validInput={validConfirmation && confirmation !== ''}>
            <FontAwesomeIcon icon={faCheck} color='#38FF5D'/>
          </RegisterValidationIcon>
          <RegisterValidationIcon validInput={!validConfirmation && confirmation !== ''}>
            <FontAwesomeIcon icon={faTimes} color='#FF6161'/>
          </RegisterValidationIcon>
        </RegisterLabel>
        <input 
          type="password"
          id='confirm'
          ref={userRef}
          autoComplete='off'
          onChange={e => setConfirmation(e.target.value)}
          required
          aria-invalid={validConfirmation ? "false" : "true"}
          aria-describedby='confirmnote'
          onFocus={() => setConfirmationFocus(true)}
          onBlur={() => setConfirmationFocus(false)}
        />
        <RegisterInfoText id='confirmnote' valid={validConfirmation && confirmation !== '' && confirmationFocus}>
          <FontAwesomeIcon icon={faInfoCircle}/>
          confirmation does not match
        </RegisterInfoText>
      </form>
      </section>
    </RegisterPageContainer>
  )
}

export default Register