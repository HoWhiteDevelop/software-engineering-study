import './App.css'
import {faker} from "@faker-js/faker"
import RestartButton from './components/RestartButton'
import Result from './components/Result'
import UserTypings from './components/UserTypings'

const words = faker.word.words(15)

function App() {

  // text-primary-500
  return (
  <div>
    <TimeLeftReminder timeLeft={20}/>
    <GenerateRandom words={words}/> 
    <UserTypings userInfo={"he llo"} className={"flex flex-row items-center text-primary-400 mt-10"}/>
    <RestartButton 
    className={"m-auto text-slate-500"} 
    onRestart = {()=>{null}}/>
    <Result 
    errors={90}
    accuracyPercentage={75}
    total={100}
    className={"mt-10"}/>
   
  </div>
)}

const GenerateRandom = ({words}:{words:string})=>{
  return(
    <div className=' text-slate-500 text-4xl'>
      words:{words}
    </div>
  )
}
const TimeLeftReminder = ({timeLeft}:{timeLeft:number})=>{
  return(
    <div className=' text-primary-500 text-xl'>
        time:{timeLeft}
    </div>
  )
}
export default App
