// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// Copyright 2024 FeatureForm Inc.
//

import { ThemeProvider } from '@mui/material/styles';
import { cleanup, render } from '@testing-library/react';
import 'jest-canvas-mock';
import React from 'react';
import HeaderBar from '../src/components/headerBar/HeaderBar';
import TEST_THEME from '../src/styles/theme';

const mockVersion = 'v0.8.1';
const apiMock = {
  fetchVersionMap: () => {
    return Promise.resolve({ data: { version: mockVersion } });
  },
};

jest.mock('next/router', () => ({
  useRouter: () => {
    return {
      push: jest.fn(),
      query: '',
    };
  },
}));

describe('HeaderBar version tests', () => {
  const VERSION_ID = 'versionPropId';

  const getTestBody = (api = apiMock) => {
    return (
      <>
        <ThemeProvider theme={TEST_THEME}>
          <HeaderBar api={api} />
        </ThemeProvider>
      </>
    );
  };

  beforeEach(() => {
    console.warn = jest.fn();
    jest.resetAllMocks();
  });

  afterEach(() => {
    cleanup();
  });

  test('Issue-812: The HeaderBar component displays the correct version data', async () => {
    //given:
    const helper = render(getTestBody());

    // when: we find the version elem
    const versionField = await helper.findByTestId(VERSION_ID);

    // then: the correct value displays
    expect(versionField).toBeDefined();
    expect(versionField.nodeName).toBe('SPAN');
    //includes mock userName separator
    expect(versionField.textContent).toBe(`Version: ${mockVersion}`);
  });
});
