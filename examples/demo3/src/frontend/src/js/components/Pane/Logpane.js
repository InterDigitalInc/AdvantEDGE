import React from 'react';
import { Typography } from '@rmwc/typography';
import { Elevation } from '@rmwc/elevation';
import Title from '@/js/components/Title';
import '@material/typography/dist/mdc.typography.css';
export default function LogPane({ data }) {
  // Generate data
  function renderActivity(data) {
    if (data) {
      return data.map((element) => {
        return (
          <Typography use="caption" style={{ display: 'block' }}>
            {element}
          </Typography>
        );
      });
    }
  }

  return (
    <div>
      <Elevation
        z={2}
        className="component-style "
        style={{ padding: 10, marginBottom: 10 }}
      >
        <Title>Activity Logs</Title>
        <div style={{ height: '45vh', overflowX: 'auto', overflowY: 'auto', whiteSpace: 'nowrap'}}>
          {renderActivity(data)}
        </div>
      </Elevation>
    </div>
  );
}
