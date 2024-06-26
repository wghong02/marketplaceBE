import React from "react";
import { Spinner } from "@chakra-ui/react";

interface LoadingWrapperProps {
  // defining the props needed for loading wrapper
  loading: boolean;
  hasFetched: boolean;
  children: React.ReactNode;
  hasData?: boolean;
}

const LoadingWrapper: React.FC<LoadingWrapperProps> = ({
  loading,
  hasFetched,
  children,
  hasData= true,
}) => {
  // wraps around a child react component to handle loading screen
  if (loading) {
    return <Spinner size="xl" />;
  }

  if (hasFetched && !hasData) {
    return <>No posts found.</>; // !!! handle when no post (maybe separate page)
  }

  return <>{children}</>;
};

export default LoadingWrapper;
