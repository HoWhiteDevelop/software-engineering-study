'use client'

import React, { useState, useEffect, useCallback } from 'react'
import { Button } from "@/components/ui/button"

const sampleTexts = [
  "The quick brown fox jumps over the lazy dog.",
  "React is a popular JavaScript library for building user interfaces.",
  "Practice makes perfect when it comes to typing.",
  "Coding is fun and rewarding, especially when you create useful applications.",
]

export function TypingGameComponent() {
  const [text, setText] = useState('')
  const [userInput, setUserInput] = useState('')
  const [startTime, setStartTime] = useState(0)
  const [wordCount, setWordCount] = useState(0)
  const [charCount, setCharCount] = useState(0)
  const [wpm, setWpm] = useState(0)
  const [accuracy, setAccuracy] = useState(100)
  const [isFinished, setIsFinished] = useState(false)

  const startGame = useCallback(() => {
    const randomIndex = Math.floor(Math.random() * sampleTexts.length)
    setText(sampleTexts[randomIndex])
    setUserInput('')
    setStartTime(Date.now())
    setWordCount(0)
    setCharCount(0)
    setWpm(0)
    setAccuracy(100)
    setIsFinished(false)
  }, [])

  useEffect(() => {
    startGame()
  }, [startGame])

  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const inputVal = e.target.value
    setUserInput(inputVal)

    // Count correctly typed characters
    let correctChars = 0
    for (let i = 0; i < inputVal.length; i++) {
      if (inputVal[i] === text[i]) {
        correctChars++
      }
    }

    // Update character count and accuracy
    setCharCount(correctChars)
    setAccuracy((correctChars / inputVal.length) * 100)

    // Check if the game is finished
    if (inputVal === text) {
      setIsFinished(true)
      const words = text.split(' ').length
      const endTime = Date.now()
      const timeInMinutes = (endTime - startTime) / 60000
      setWordCount(words)
      setWpm(Math.round(words / timeInMinutes))
    }
  }

  return (
    <div className="max-w-2xl mx-auto p-4 space-y-4">
      <h1 className="text-2xl font-bold text-center">打字练习游戏</h1>
      <div className="bg-secondary p-4 rounded">
        {text.split('').map((char, index) => (
          <span
            key={index}
            className={
              index < userInput.length
                ? char === userInput[index]
                  ? 'text-green-500'
                  : 'text-red-500'
                : ''
            }
          >
            {char}
          </span>
        ))}
      </div>
      <input
        type="text"
        value={userInput}
        onChange={handleInputChange}
        className="w-full p-2 border rounded"
        placeholder="开始输入..."
        disabled={isFinished}
      />
      <div className="flex justify-between text-sm">
        <span>准确率: {accuracy.toFixed(2)}%</span>
        <span>WPM: {wpm}</span>
        <span>字符数: {charCount}</span>
      </div>
      <Button onClick={startGame} className="w-full">
        {isFinished ? '再来一次' : '重新开始'}
      </Button>
    </div>
  )
}
