import Caret from "./Caret";

const UserTypings = ({
  userInfo,
  className,
}: {
  userInfo: string;
  className?: string;
}) => {
  //此正则替换的方法亦会导致更严重的不兼容，随机的英文字词被割裂
  //const typeCharacter = userInfo.replace(/ /g, "\u00A0").split("");
  const typeCharacter = userInfo.split('');
  console.log(typeCharacter);
  return (
    <div className={className}>
      {typeCharacter.map((item, index) => (
        <span key={`${index}_${item}`}>{item}</span>
      ))}
      <Caret></Caret>
    </div>
  );
};

export default UserTypings;
