import useWords from "./useWords";
import { useState } from "react";

export type State = "start" | "run" | "finish";

const useEngine = () => {
  const [state, setState] = useState<State>("start");

  const { words, updateWords } = useWords(6);

  return { state, words, updateWords };
};

export default useEngine;
