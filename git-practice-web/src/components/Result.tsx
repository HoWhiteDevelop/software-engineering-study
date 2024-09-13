import {motion} from "framer-motion"
import { formatPercentage } from "../utils/helpers";

interface ModuleResult{
    errors?:number,
    accuracyPercentage?:number,
    total?:number,
    className?:string
}

const Result = ({
    errors,
    accuracyPercentage,
    total,
    className
}:ModuleResult)=>{
    const initial = {opacity: 0};
    const animate = {opacity: 1};
    const duration = {duration:6};

    return (
        <motion.ul className={`${className} flex flex-col items-center text-primary-400 space-y-3`}>
           <motion.li
            initial={initial}
            animate={animate}
            transition={{...duration, delay:0}}
           >Result</motion.li>
           <motion.li
           initial={initial}
           animate={animate}
           transition={{...duration, delay:0.5}}
           >Accuracy:{formatPercentage(accuracyPercentage)}</motion.li>
           <motion.li
           initial={initial}
           animate={animate}
           transition={{...duration, delay:1.0}}
           className=" text-red-500"
           >errors:{errors}</motion.li>
           <motion.li
           initial={initial}
           animate={animate}
           transition={{...duration, delay:1.3}}
           >total:{total}</motion.li>
        </motion.ul>
    )
}
export default Result