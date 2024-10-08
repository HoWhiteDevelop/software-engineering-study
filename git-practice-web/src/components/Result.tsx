import { motion } from "framer-motion";
import { formatPercentage } from "../utils/helpers";

interface ModuleResult {
  errors: number;
  accuracyPercentage: number;
  total: number;
  className: string;
}

const Result = ({
  errors,
  accuracyPercentage,
  total,
  className,
}: ModuleResult) => {
  const initial = { opacity: 0 };
  const animate = { opacity: 1 };
  const duration = { duration: 6 };

  return (
    <motion.ul
      className={`${className} flex flex-col items-center text-primary-400 space-y-4`}
    >
      <motion.li
        initial={initial}
        animate={animate}
        //必须...duration展开对象,否则delay无法与duration:6处在一个对象体之中
        transition={{ ...duration, delay: 0 }}
      >
        Result
      </motion.li>
      <motion.li
        initial={initial}
        animate={animate}
        transition={{ ...duration, delay: 0.5 }}
      >
        Accuracy:{formatPercentage(accuracyPercentage)}
      </motion.li>
      <motion.li
        initial={initial}
        animate={animate}
        transition={{ ...duration, delay: 1.0 }}
        className=" text-red-600 brightness-200"
      >
        Errors:{errors}
      </motion.li>
      <motion.li
        initial={initial}
        animate={animate}
        transition={{ ...duration, delay: 1.4 }}
      >
        Total:{total}
      </motion.li>
    </motion.ul>
  );
};
export default Result;
