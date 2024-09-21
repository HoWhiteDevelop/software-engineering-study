import { motion } from "framer-motion";

const Caret = () => {
  return (
    <motion.div
      className=" inline-block w-0.5 h-7 bg-primary-500"
      initial={{ opacity: 0 }}
      animate={{ opacity: 1 }}
      transition={{ repeat: Infinity, duration: 0.5, ease: "easeInOut" }}
      exit={{opacity:1}}
    ></motion.div>
  );
};

export default Caret;
