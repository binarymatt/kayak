import { useParams } from "@solidjs/router";

export const Stream = () => {
  const params = useParams();
  const name = params.name;
  return <>{name} stream detail coming soon</>;
};
