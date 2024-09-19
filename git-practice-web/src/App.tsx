import "./App.css";
import { faker } from "@faker-js/faker";
import RestartButton from "./components/RestartButton";
import Result from "./components/Result";
import UserTypings from "./components/UserTypings";

const words = faker.word.words(15);

function App() {
  // text-primary-500
  return (
    <div>
      <TimeLeftReminder timeLeft={20} />
      <div className=" relative text-3xl max-w-xl leading-relaxed break-all">
        <GenerateRandom words={words} />
        <UserTypings
          userInfo={words}
          className={" absolute inset-0 text-red-500"}
        />
      </div>
      <RestartButton
        className={"m-auto text-slate-500"}
        onRestart={() => {
          null;
        }}
      />
      <Result
        errors={90}
        accuracyPercentage={75}
        total={100}
        className={"mt-10"}
      />
    </div>
  );
}

const GenerateRandom = ({ words }: { words: string }) => {
  return <div className=" text-slate-400">{words}</div>;
};
const TimeLeftReminder = ({ timeLeft }: { timeLeft: number }) => {
  return <div className=" text-primary-500 text-xl">time:{timeLeft}</div>;
};
export default App;
