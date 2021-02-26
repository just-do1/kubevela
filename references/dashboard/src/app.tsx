import React from 'react';

import { notification } from 'antd';
import { RequestConfig, RunTimeLayoutConfig } from 'umi';
import { ResponseError } from 'umi-request';

import Footer from '@/components/Footer';
import RightContent from '@/components/RightContent';
import { PageLoading, Settings as LayoutSettings } from '@ant-design/pro-layout';

import defaultSettings from '../config/defaultSettings';

export const initialStateConfig = {
  loading: <PageLoading />,
};

export async function getInitialState(): Promise<{
  settings?: LayoutSettings;
}> {
  return {
    settings: defaultSettings,
  };
}

export const layout: RunTimeLayoutConfig = ({ initialState }) => {
  return {
    rightContentRender: () => <RightContent />,
    disableContentMargin: false,
    footerRender: () => <Footer />,
    menuHeaderRender: undefined,
    // custom 403 page
    // unAccessible: <div>unAccessible</div>,
    ...initialState?.settings,
  };
};

const codeMessage = {
  200: 'The server successfully returned the requested data. ',
  201: 'Create or modify data successfully. ',
  202: 'A request has entered the background queue (asynchronous task). ',
  204: 'Delete data successfully. ',
  400: 'There was an error in the request sent, and the server did not create or modify data. ',
  401: 'The user does not have permission (the token, username, password are wrong). ',
  403: 'The user is authorized, but access is forbidden. ',
  404: 'The request sent is for a record that does not exist, and the server is not operating. ',
  405: 'The requested method is not allowed. ',
  406: 'The requested format is not available. ',
  410: 'The requested resource has been permanently deleted and will no longer be available. ',
  422: 'When creating an object, a validation error occurred. ',
  500: 'An error occurred in the server, please check the server. ',
  502: 'Gateway error. ',
  503: 'The service is unavailable, the server is temporarily overloaded or maintained. ',
  504: 'The gateway has timed out. ',
};

const errorHandler = (error: ResponseError) => {
  const { response, data } = error;
  if (response && response.status) {
    const errorText = data?.data || codeMessage[response.status] || response.statusText;
    const { status, url } = response;

    notification.error({
      message: `Request error ${status}: ${url}`,
      description: errorText,
    });
  }

  if (!response) {
    notification.error({
      description: 'Your network is abnormal and cannot connect to the server',
      message: 'Network failure',
    });
  }
  throw error;
};

export const request: RequestConfig = {
  errorHandler,
};
