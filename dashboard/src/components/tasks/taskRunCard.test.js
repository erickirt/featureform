// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// Copyright 2024 FeatureForm Inc.
//

import { ThemeProvider } from '@mui/material/styles';
import { cleanup, fireEvent, render } from '@testing-library/react';
import 'jest-canvas-mock';
import React from 'react';
import TEST_THEME from '../../styles/theme';
import TaskRunCard from './taskRunCard';
import { taskCardDetailsResponse, taskRunsResponse } from './test_data';

const dataAPIMock = {
  getTaskRuns: jest.fn().mockResolvedValue(taskRunsResponse),
  getTaskRunDetails: jest.fn().mockResolvedValue(taskCardDetailsResponse),
};

jest.mock('../../hooks/dataAPI', () => ({
  useDataAPI: () => {
    return dataAPIMock;
  },
}));

describe('Task run card detail tests', () => {
  const REFRESH_ICON_ID = 'taskRunRefreshIcon';
  const getTestBody = (taskId, taskRunId) => {
    return (
      <>
        <ThemeProvider theme={TEST_THEME}>
          <TaskRunCard
            handleClose={jest.fn()}
            taskId={taskId}
            taskRunId={taskRunId}
          />
        </ThemeProvider>
      </>
    );
  };

  beforeEach(() => {});

  afterEach(() => {
    jest.resetAllMocks();
    jest.restoreAllMocks();
    cleanup();
  });

  test('Basic task card detail render with refresh', async () => {
    //given:
    jest.useFakeTimers();
    const taskRunId = taskCardDetailsResponse.taskRun.runId;
    const taskId = taskCardDetailsResponse.taskRun.taskId;
    const helper = render(getTestBody(taskId, taskRunId));

    const foundRefreshBtn = await helper.findByTestId(REFRESH_ICON_ID);
    fireEvent.click(foundRefreshBtn);
    jest.advanceTimersByTime(1500);
    await helper.findByTestId(REFRESH_ICON_ID);

    //expect:
    expect(dataAPIMock.getTaskRunDetails).toHaveBeenCalledTimes(2);
    expect(dataAPIMock.getTaskRunDetails).toHaveBeenCalledWith(
      taskId,
      taskRunId
    );
  });
});
