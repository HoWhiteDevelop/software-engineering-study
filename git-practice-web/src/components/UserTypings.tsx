const UserTypings = ({
  userInfo,
  className,
}: {
  userInfo: string;
  className?: string;
}) => {
  const typeCharacter = userInfo.replace(/ /g, "\u00A0").split("");
  console.log(typeCharacter);
  return (
    <div className={className}>
      {typeCharacter.map((item, index) => (
        <span key={index}>{item}</span>
      ))}
    </div>
  );
};

export default UserTypings;
