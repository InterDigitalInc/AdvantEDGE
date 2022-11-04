/*
 * Copyright (c) 2022  The AdvantEDGE Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

import React from 'react';
import { useState, useEffect } from 'react';

// material ui
import { Grid, GridCell } from '@rmwc/grid';
import '@material/elevation/dist/mdc.elevation.css';
import '@material/layout-grid/dist/mdc.layout-grid.css';
import Button from '@mui/material/Button';
import '@material/button/dist/mdc.button.css';
import {
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  DialogButton
} from '@rmwc/dialog';
import '@material/dialog/dist/mdc.dialog.css';
import '@material/textfield/dist/mdc.textfield.css';
import '@material/floating-label/dist/mdc.floating-label.css';
import '@material/notched-outline/dist/mdc.notched-outline.css';
import '@material/line-ripple/dist/mdc.line-ripple.css';
import { TextField } from '@rmwc/textfield';
import Container from '@mui/material/Container';
import { createTheme, ThemeProvider } from '@mui/material/styles';

// custom package
import Header from '@/js/components/Layout/Header';
import ActivityPane from '@/js/components/Pane/AppInfo';
import LogPane from '@/js/components/Pane/Logpane';
import AmsPane from '@/js/components/Pane/Amspane';
const { palette } = createTheme();
// Import client
import * as demoSvcRestApiClient from '../../../../client/src/index.js';

// import css
import '@/css/global.css';

export default function Homepage() {
  // MEEP Demo REST API JS Client
  // Configure server url based on environmental variable externally or preset internally

  var basepath;

  basepath = 'http://' + location.host + location.pathname;


  demoSvcRestApiClient.ApiClient.instance.basePath = basepath.replace(
    /\/+$/,
    ''
  );
  var appInfoApi = new demoSvcRestApiClient.FrontendApi();

  // State management
  const [appInfo, setAppInfo] = useState({});
  const [appLog, setAppLog] = useState([]);
  const [amsLog, setAmsLog] = useState([]);
  const [textValue, setTextValue] = useState('');

  // Button
  const [registered, setRegisteration] = useState(false);
  const [start, setStart] = useState(false);
  const [modal, setModal] = useState(false);
  const [amsModal, setAmsModal] = useState(false);

  // Inital loading
  useEffect(() => {
    appInfoApi.getPlatformInfo((error, data, response) => {
      if (error !== null) {
        // console.log(error);
      } else {
        setAppInfo(response.body);
      }
    });
  }, []);


  // If app is registered & added terminal device
  // Perform polling on ams
  useEffect(() => {
    if (registered) {

      const interval = setInterval(() => {

        appInfoApi.getAmsDevices((error, data, response) => {
          if (error !== null) {
            // console.log(error);
          } else {
            setAmsLog(response.body);
          }
        });
      }, 1000);
      return () => clearInterval(interval);
    }
  }, [amsLog, registered]);

  // If app is registered or app info changes
  // Perform polling on app info
  useEffect(() => {
    if (registered) {

      const interval = setInterval(() => {
        appInfoApi.getPlatformInfo((error, data, response) => {
          if (error !== null) {
            //      console.log(error);
          } else {
            setAppInfo(response.body);
          }
        });
      }, 1000);
      return () => clearInterval(interval);
    }
  }, [appInfo, registered]);

  // If app is registered or activity log changes
  // Peform polling activity logs
  useEffect(() => {
    if (registered) {
      const interval = setInterval(() => {
        appInfoApi.getActivityLogs((error, data, response) => {
          if (error !== null) {
            //console.log(error);
          } else {
            setAppLog(response.body);
          }
        });
      }, 1000);
      return () => clearInterval(interval);
    }
  }, [appLog, registered]);

  // Stop polling app info
  // Turn registeration to false + clean activity info + one request to app info + one request to ams
  const deRegisterapp = async () => {
    setAppInfo({});
    setAmsLog([]);
    setRegisteration(false);
    appInfoApi.deregister((err, data, resp) => {
      if (err !== null) {
        // console.log(err);
      }
      appInfoApi.getActivityLogs((error, data, response) => {
        if (error !== null) {
        //  console.log(error);
        } else {
          setAppLog(response.body);
        }
      });
      appInfoApi.getPlatformInfo((error, data, response) => {
        if (error !==null) {
          // console.log(error);
        } else {
          setAppInfo(response.body);
        }
      });
    });
  };


  const removeAmsDevice = async (device) => {
    appInfoApi.deleteAmsDevice(device, (err, data, resp) => {
      if (err !== null) {
        // console.log(err);
      }
    });
  };

  const addTerminalDevices = () => {
    appInfoApi.updateAmsDevices(textValue, (err, data, resp) => {
      if (err !== null) {
        // console.log(err);
      }
    });
  }
  ;

  return (
    <ThemeProvider theme={theme}>
      <div>
        <Header></Header>
        <Container maxWidth="100vw" sx={{ mt: 2 }}>
          <Grid style={styles.headlineGrid}>
            <GridCell span="5">
              <ActivityPane data={appInfo}></ActivityPane>
            </GridCell>
            <GridCell span="2">
              <div>
                <Button
                  outlined
                  color="themecolor"
                  disabled={registered}
                  fullWidth
                  variant="contained"
                  sx={{ mt: 3, mb: 2 }}
                  onClick={() => {
                    appInfoApi.register((err, data, resp) => {
                      if (err !== null) {
                        // console.log(err);
                      } else {
                        setAppInfo(resp.body);
                        setRegisteration(true);
                      }
                    });
                  }}
                >
                  Register Application
                </Button>
              </div>
              <div>
                <Button
                  outlined
                  color="themecolor"
                  disabled={!registered}
                  fullWidth
                  variant="contained"
                  sx={{ mt: 3, mb: 2 }}
                  onClick={() => {
                    setStart(true);
                  }}
                >
                  De-Register Application
                </Button>

                <Dialog
                  open={start}
                  onClose={() => {
                    setStart(false);
                  }}
                >
                  <DialogTitle theme="primary" style={styles.title}>
                    Confirm
                  </DialogTitle>
                  <DialogContent>
                    Clear Mec Resource for Mec Application Demo 3?
                  </DialogContent>
                  <DialogActions>
                    <DialogButton action="close">Cancel</DialogButton>
                    <DialogButton
                      action="accept"
                      isDefaultAction
                      onClick={deRegisterapp}
                    >
                      Confirm
                    </DialogButton>
                  </DialogActions>
                </Dialog>
              </div>

              <div style={{ marginTop: '5rem' }}>
                <Button
                  outlined
                  color="themecolor"
                  disabled={!registered}
                  fullWidth
                  variant="contained"
                  sx={{ mt: 3, mb: 2 }}
                  onClick={() => {
                    setModal(true);
                  }}
                >
                  Add AMS Device
                </Button>
              </div>
              <Dialog
                open={modal}
                onClose={() => {
                  setModal(false);
                }}
              >
                <DialogContent>Update mobility service resource</DialogContent>
                <TextField
                  style={{ margin: '10px' }}
                  label="device"
                  value={textValue}
                  onChange={(event) => {
                    setTextValue(event.target.value);
                  }}
                ></TextField>
                <DialogActions>
                  <DialogButton action="close">Cancel</DialogButton>
                  <DialogButton
                    action="accept"
                    isDefaultAction
                    onClick={() => {
                      addTerminalDevices();
                    }}
                  >
                    Confirm
                  </DialogButton>
                </DialogActions>
              </Dialog>
              <Button
                outlined
                color="themecolor"
                disabled={!registered}
                fullWidth
                variant="contained"
                sx={{ mt: 3, mb: 2 }}
                onClick={() => {
                  setAmsModal(true);
                }}
              >
                Remove AMS Device
              </Button>

              <Dialog
                open={amsModal}
                onClose={() => {
                  setAmsModal(false);
                }}
              >
                <DialogTitle theme="primary" style={styles.title}>
                  Confirm
                </DialogTitle>
                <DialogContent>Delete device from AMS Resource?</DialogContent>
                <TextField
                  style={{ margin: '10px' }}
                  label="device"
                  value={textValue}
                  onChange={(event) => {
                    setTextValue(event.target.value);
                  }}
                ></TextField>
                <DialogActions>
                  <DialogButton action="close">Cancel</DialogButton>
                  <DialogButton
                    action="accept"
                    isDefaultAction
                    onClick={() => removeAmsDevice(textValue)}
                  >
                    Confirm
                  </DialogButton>
                </DialogActions>
              </Dialog>
            </GridCell>
            <GridCell span="5" className={{ height: '100vh' }}>
              <LogPane data={appLog}></LogPane>
              <AmsPane data={amsLog}></AmsPane>
            </GridCell>
          </Grid>
        </Container>
      </div>
    </ThemeProvider>
  );
}

const styles = {
  headlineGrid: {
    marginBottom: 10
  },
  title: {
    paddingTop: 10,
    paddingBottom: 15
  }
};

const theme = createTheme({
  palette: {
    themecolor: palette.augmentColor({ color: { main: '#379DD8' } })
  }
});
