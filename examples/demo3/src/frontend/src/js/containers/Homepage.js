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
 
  // var basepath = process.env.URL;
  var basepath = 'http://' + location.host + location.pathname;
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
  const [amsStart, setAmsStart] = useState(false);
  const [amsModal, setAmsModal] = useState(false);

  // Inital loading 
  useEffect(() => {
    appInfoApi.infoApplicationMecPlatformGet((error, data, response) => {
      if (error !== null) {
        console.log(error);
      } else {
        setAppInfo(response.body);
      }
    });
  }, []);

  // If app is registered & added terminal device
  // Perform polling on ams 
  useEffect(() => {
    if (registered && amsStart) {
      const interval = setInterval(() => {
        appInfoApi.infoAmsLogsGet(20000, (error, data, response) => {
          if (error !== null) {
            console.log(error);
          } else {
            setAmsLog(response.body);
          }
        });
      }, 1000);
      return () => clearInterval(interval);
    }
  }, [amsLog, amsStart, registered]);

  // If app is registered or app info changes
  // Perform polling on app info
  useEffect(() => {
    if (registered) {
      const interval = setInterval(() => {
        appInfoApi.infoApplicationMecPlatformGet((error, data, response) => {
          if (error != null) {
            console.log(error);
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
        appInfoApi.infoLogsGet(2000, (error, data, response) => {
          if (error != null) {
            console.log(error);
          } else {
            setAppLog(response.body);
            console.log(response.body);
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
    setTerminalDevices([]);
    appInfoApi.infoApplicationMecPlatformDeleteDelete((err, data, resp) => {
      if (err != null) {
        console.log(err);
      }
      appInfoApi.infoLogsGet(40, (error, data, response) => {
        if (error != null) {
          console.log(error);
        } else {
          setAppLog(response.body);
        }
      });
      appInfoApi.infoApplicationMecPlatformGet((error, data, response) => {
        if (error != null) {
          console.log(error);
        } else {
          setAppInfo(response.body);
        }
      });
    });
  };

  const [terminalDevices, setTerminalDevices] = useState([]);

  const removeAmsDevice = async (device) => {
    let mutableArray = terminalDevices;
    const result = mutableArray.filter(e => e != device );
    appInfoApi.serviceAmsDeleteDeviceDelete(device, (err, data, resp) => {
      if (err != null) {
        console.log(err);
      }
      setTerminalDevices(result);
    });
  };

  const addTerminalDevices = () => {
    let devices = terminalDevices;
    let deviceMap = new Set(devices);
    if (!deviceMap.has(textValue)) {
      appInfoApi.serviceAmsUpdateDevicePut(textValue, (err, data, resp) => {
        if (err != null) {
          console.log(err);
        }
        setTerminalDevices(e => [... e, textValue]); 
      });
    }      
  };

  return (
    <ThemeProvider theme={theme}>
      <div className="ui-background" style={{ height: '100%' }}>
        <Header></Header>
        <Container maxWidth="100vw" sx={{ mt: 2, mb: 4 }}>
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
                    appInfoApi.registerAppMecPlatformPost(
                      (err, data, resp) => {
                        if (err != null) {
                          console.log(err);
                        } else {
                          setAppInfo(resp.body);
                          setRegisteration(true);
                        }
                      }
                    );
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
                    onClick={(e) => {
                      setAmsStart(true);
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
            <GridCell span="5">
              <LogPane data={appLog}></LogPane>
              <AmsPane data={amsLog} style={{ marginTop: '5px' }}></AmsPane>
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
    themecolor: palette.augmentColor({ color: { main: '#379DD8'} })
  }
});
