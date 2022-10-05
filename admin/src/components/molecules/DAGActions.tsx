import { Stack } from '@mui/material';
import React from 'react';
import { SchedulerStatus, Status } from '../../models';
import ActionButton from '../atoms/ActionButton';
import { useNavigate } from 'react-router-dom';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faPlay, faStop, faReply } from '@fortawesome/free-solid-svg-icons';
import VisuallyHidden from '../atoms/VisuallyHidden';

type LabelProps = {
  show: boolean;
  children: React.ReactNode;
};

type Props = {
  status?: Status;
  name: string;
  label?: boolean;
  redirectTo?: string;
  refresh?: () => void;
};

function Label({ show, children }: LabelProps): JSX.Element {
  if (show) return <>{children}</>;
  return <VisuallyHidden>{children}</VisuallyHidden>;
}

function DAGActions({
  status,
  name,
  refresh,
  redirectTo,
  label = true,
}: Props) {
  const nav = useNavigate();

  const onSubmit = React.useCallback(
    async (
      warn: string,
      params: {
        name: string;
        action: string;
        requestId?: string;
      }
    ) => {
      const form = new FormData();
      if (params.action == 'start') {
        const parameters = window.prompt(
          'Enter parameters (for default parameters, leave blank and click OK).',
          ''
        );
        if (parameters === null) {
          //hint cancel
          return;
        }
        form.set('params', parameters);
      } else {
        if (!confirm(warn)) {
          return;
        }
      }
      form.set('action', params.action);
      if (params.requestId) {
        form.set('request-id', params.requestId);
      }
      const url = `${API_URL}/dags/${params.name}`;
      const ret = await fetch(url, {
        method: 'POST',
        mode: 'cors',
        body: form,
      });
      if (redirectTo) {
        nav(redirectTo);
        refresh && refresh();
        return;
      }
      if (!ret.ok) {
        const e = await ret.text();
        alert(e || 'Failed to submit');
      }
      refresh && refresh();
    },
    [refresh]
  );
  const buttonState = React.useMemo(
    () => ({
      start: status?.Status != SchedulerStatus.Running,
      stop: status?.Status == SchedulerStatus.Running,
      retry:
        status?.Status != SchedulerStatus.Running && status?.RequestId != '',
    }),
    [status]
  );
  return (
    <Stack direction="row" spacing={2}>
      <ActionButton
        label={label}
        icon={
          <>
            <Label show={false}>Start</Label>
            <span className="icon">
              <FontAwesomeIcon icon={faPlay} />
            </span>
          </>
        }
        disabled={!buttonState['start']}
        onClick={() =>
          onSubmit('Do you really want to start the DAG?', {
            name: name,
            action: 'start',
          })
        }
      >
        {label && 'Start'}
      </ActionButton>
      <ActionButton
        label={label}
        icon={
          <>
            <Label show={false}>Stop</Label>
            <span className="icon">
              <FontAwesomeIcon icon={faStop} />
            </span>
          </>
        }
        disabled={!buttonState['stop']}
        onClick={() =>
          onSubmit('Do you really want to cancel the DAG?', {
            name: name,
            action: 'stop',
          })
        }
      >
        {label && 'Stop'}
      </ActionButton>
      <ActionButton
        label={label}
        icon={
          <>
            <Label show={false}>Retry</Label>
            <span className="icon">
              <FontAwesomeIcon icon={faReply} />
            </span>
          </>
        }
        disabled={!buttonState['retry']}
        onClick={() =>
          onSubmit(
            `Do you really want to rerun the last execution (${status?.RequestId}) ?`,
            {
              name: name,
              requestId: status?.RequestId,
              action: 'retry',
            }
          )
        }
      >
        {label && 'Retry'}
      </ActionButton>
    </Stack>
  );
}
export default DAGActions;
