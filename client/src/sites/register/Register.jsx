import React from 'react'
import { RegisterLabel, RegisterPageContainer, RegisterValidationIcon } from './styles'
import { useRef, useState, useEffect } from 'react'
import {faCheck, faTimes, faInfoCircle} from '@fortawesome/free-solid-svg-icons'
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome'

function Register() {
  const EMAIL_REGEX = /^[A-Za-z0-9._%+-§$%!?]+@[A-Za-z0-9.-]+\.[a-z]{2,4}$/
  const PASSWORD_REGEX = /^(?=.*[a-z])(?=.*[A-Z])(?=.*[0-9])(?=.*[!?&%$§#@€]).{10,}$/

  const userRef = useRef()
  const errRef = useRef()

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
          <RegisterValidationIcon validInput={validEmail && user != ''}>
            <FontAwesomeIcon icon={faCheck}/>
          </RegisterValidationIcon>
          <RegisterValidationIcon validInput={!validEmail && user != ''}>
            <FontAwesomeIcon icon={faTimes}/>
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
        <p id='uidnote'>
          <FontAwesomeIcon icon={faInfoCircle}/>
          invalid email
        </p>

        <RegisterLabel htmlFor="password">
          Email:
          <RegisterValidationIcon validInput={validPassword && password != ''}>
            <FontAwesomeIcon icon={faCheck}/>
          </RegisterValidationIcon>
          <RegisterValidationIcon validInput={!validPassword && password != ''}>
            <FontAwesomeIcon icon={faTimes}/>
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
        <p id='pwdnote'>
          <FontAwesomeIcon icon={faInfoCircle}/>
          password must have at least 10 characters and lowercase and uppercase letters. also it needs to have digits and special characters like %&$§!?#€&
        </p>

        <RegisterLabel htmlFor="confirm">
          Email:
          <RegisterValidationIcon validInput={validConfirmation && confirmation != ''}>
            <FontAwesomeIcon icon={faCheck}/>
          </RegisterValidationIcon>
          <RegisterValidationIcon validInput={!validConfirmation && confirmation != ''}>
            <FontAwesomeIcon icon={faTimes}/>
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
        <p id='confirmnote'>
          <FontAwesomeIcon icon={faInfoCircle}/>
          confirmation does not match
        </p>
      </form>
      </section>
    </RegisterPageContainer>
  )
}

export default Register