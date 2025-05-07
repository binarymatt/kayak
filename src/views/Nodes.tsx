import { useParams } from "@solidjs/router";
export const Nodes = () => {
  const params = useParams();
  console.log(params.index);
  return <>Nodes soon - {params.index}</>;
};
