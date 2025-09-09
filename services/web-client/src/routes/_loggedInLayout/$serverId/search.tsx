import { createFileRoute } from "@tanstack/react-router";

const SearchLayout = () => {
  return <div>Search page</div>;
};

export const Route = createFileRoute("/_loggedInLayout/$serverId/search")({
  component: SearchLayout,
});
