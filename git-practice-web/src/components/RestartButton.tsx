import { useRef } from "react"
import {MdRefresh} from "react-icons/md"

interface RestartButtonProps{
    onRestart:()=>void
    className?:string
}

const RestartButton = ({className,onRestart:handleRestart}:RestartButtonProps)=>{
    const handleClick = ()=>{
        buttonRef.current?.blur()
        handleRestart()
    }//函数需额外声明引用
    const buttonRef = useRef<HTMLButtonElement>(null)//绑定ref，消除后续点击焦点无法自动移除的问题

    return(
        <button
        ref = {buttonRef}
        onClick={handleClick} 
        className={`block rounded px-8 py-3 hover:bg-slate-700/30 my-4 ${className}`}>
            <MdRefresh className=" w-6 h-6"/>
        </button>
    )
}

export default RestartButton