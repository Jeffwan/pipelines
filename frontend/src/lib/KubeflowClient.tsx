import React, { useLayoutEffect, useContext, useRef, FC, useState, useEffect } from 'react';
import { logger } from './Utils';

declare global {
  interface Window {
    // Provided by:
    // 1. https://github.com/kubeflow/kubeflow/tree/master/components/centraldashboard#client-side-library
    // 2. /frontend/server/server.ts -> KUBEFLOW_CLIENT_PLACEHOLDER
    centraldashboard: any;
  }
}

let namespace: string | undefined;
let registeredHandler: undefined | ((namespace: string) => void);
function onNamespaceChanged(handler: (namespace: string) => void) {
  registeredHandler = handler;
}

export function init(): void {
  try {
    // Init method will invoke the callback with the event handler instance
    // and a boolean indicating whether the page is iframed or not
    window.centraldashboard.CentralDashboardEventHandler.init((cdeh: any) => {
      // Binds a callback that gets invoked anytime the Dashboard's
      // namespace is changed
      cdeh.onNamespaceSelected = (newNamespace: string) => {
        namespace = newNamespace;
        if (registeredHandler) {
          registeredHandler(namespace);
        }
      };
    });
  } catch (err) {
    logger.error('Failed to initialize central dashboard client', err);
  }
}

export const NamespaceContext = React.createContext<string | undefined>(undefined);
export class NamespaceContextProvider extends React.Component {
  state = {
    namespace,
  };
  componentDidMount() {
    onNamespaceChanged(ns => this.setState({ namespace: ns }));
  }
  render() {
    return <NamespaceContext.Provider value={this.state.namespace} {...this.props} />;
  }
}

export interface EurusMetadata {
  username?: string;
  jwtToken?: string;
}


function parseJwt(token: string): string {
  try {
    const base64Url =  token.split('.')[1];
    const base64 = base64Url.replace(/-/g, '+').replace(/_/g, '/');
    return  JSON.parse(decodeURIComponent(escape(atob(base64))));
  } catch (err) {
    console.error('jwt parse error', err);
    return  "";
  }
};

export const EurusContext = React.createContext<EurusMetadata>({});
export const EurusContextProvider: FC<{}> = props => {
  const [metadata, setMetadata] = useState<EurusMetadata>({});
  useEffect(() => {
    // since React is hard to retrieve the http header, we determine to get query param in frontend side.
    const queryParamsString = window.location.search.substr(1)
    if (queryParamsString === undefined || queryParamsString.length === 0) {
      console.warn('Can not find query param, skip getting user claim')
      return
    }
    const urlParams = new URLSearchParams(queryParamsString);
    const jwtToken = urlParams.get('token')
    const claims = parseJwt(jwtToken || "")
    const username = claims['username']
    const obj: EurusMetadata = { username: username!, jwtToken: jwtToken!};
    setMetadata(obj)

    console.log('---------useEffect------------')
    console.log(jwtToken)
    console.log(username)
    console.log(obj)
    console.log(metadata)
  },[])
  return <EurusContext.Provider value={metadata} {...props}></EurusContext.Provider>
}

function usePrevious<T>(value: T) {
  const ref = useRef(value);
  useLayoutEffect(() => {
    ref.current = value;
  });
  return ref.current;
}

export function useNamespaceChangeEvent(): boolean {
  const currentNamespace = useContext(NamespaceContext);
  const previousNamespace = usePrevious(currentNamespace);

  if (!previousNamespace) {
    // Previous namespace hasn't been initialized, this does not count as a change.
    // When the webapp inits, the first render will have namespace=undefined, so
    // this situation happens often.
    return false;
  }

  return previousNamespace !== currentNamespace;
}
