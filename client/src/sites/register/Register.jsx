import React from 'react'
import { RegisterPageContainer } from './styles'
import { useRef, useState, useEffect } from 'react'
import {faCheck, faTimes, faInfoCircle} from '@fortawesome/free-solid-svg-icons'
import { FontAwesomeICON } from '@fortawesome/react-fontawesome'

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
      Register
    </RegisterPageContainer>
  )
}

export default Register