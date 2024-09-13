const UserTypings = ({
    userInfo,
    className
}:{
    userInfo:string,
    className?:string
})=>{
    const typeCharacter = userInfo.split(' ')
    console.log(typeCharacter,1)
    return(
        <div className={`${className} bg-slate-600`}>
            Hello HoWhite 
            <br />
            Hello World 
        </div>
    )
}

export default UserTypings