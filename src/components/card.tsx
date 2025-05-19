import { ParentProps } from "solid-js";
const Card = (props: ParentProps) => {
  return (
    <div class="card bg-base-100 w-96 shadow-sm">
      <div class="card-body">{props.children}</div>
    </div>
  );
};

export default Card;
