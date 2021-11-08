import React from "react";
import { useState, useEffect } from "react";

// material ui
import { Grid, GridCell } from "@rmwc/grid";
import "@material/elevation/dist/mdc.elevation.css";
import "@material/layout-grid/dist/mdc.layout-grid.css";
import Button from "@mui/material/Button";
import "@material/button/dist/mdc.button.css";
import {
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  DialogButton,
} from "@rmwc/dialog";
import "@material/dialog/dist/mdc.dialog.css";
import "@material/textfield/dist/mdc.textfield.css";
import "@material/floating-label/dist/mdc.floating-label.css";
import "@material/notched-outline/dist/mdc.notched-outline.css";
import "@material/line-ripple/dist/mdc.line-ripple.css";
import { TextField } from "@rmwc/textfield";
import Container from "@mui/material/Container";

// custom package
import Header from "@/js/components/Layout/Header";
import ActivityPane from "@/js/components/Pane/AppInfo";
import LogPane from "@/js/components/Pane/logpane";
import AmsPane from "@/js/components/Pane/amspane";

// Import client
import * as demoSvcRestApiClient from "../../../../client/src/index.js";

// import css
import "@/css/global.css";

export default function Homepage() {
  // MEEP Demo REST API JS Client & configure server url
  var basepath = "http://" + location.host + location.pathname;
  // var subStr1 = basepath.split(":");
  // var subStr2 = subStr1[2].split("/");
  // var portApp = subStr2[0];
  // basepath = basepath.replace(portApp, "8093");
  // demoSvcRestApiClient.ApiClient.instance.basePath = process.env.API_URL;
  demoSvcRestApiClient.ApiClient.instance.basepath = basepath;
  var appInfoApi = new demoSvcRestApiClient.FrontendApi();

  // State management
  const [appInfo, setAppInfo] = useState({});
  const [appLog, setAppLog] = useState([]);
  const [amsLog, setAmsLog] = useState([]);
  const [textValue, setTextValue] = useState("");

  // Button
  const [registered, setRegisteration] = useState(false);
  const [start, setStart] = useState(false);
  const [modal, setModal] = useState(false);
  const [amsStart, setAmsStart] = useState(false);
  const [amsModal, setAmsModal] = useState(false);

  // Inital loading client request
  useEffect(() => {
    appInfoApi.infoApplicationMecPlatformGet(1, (error, data, response) => {
      if (error != null) {
        console.log(error);
      } else {
        setAppInfo(response.body);
        console.log(appInfo);
      }
    });
  }, []);

  // Perform polling ams log if registered app & ams button pressed
  useEffect(() => {
    if (registered && amsStart) {
      const interval = setInterval(() => {
        appInfoApi.infoAmsLogsGet(10, (error, data, response) => {
          if (error != null) {
            console.log(error);
          } else {
            setAmsLog(response.body);
          }
        });
      }, 1000);
      return () => clearInterval(interval);
    }
  }, [amsLog, amsStart, registered]);

  // Pefrom polling app info only if registered or app info changes
  useEffect(() => {
    if (registered) {
      const interval = setInterval(() => {
        appInfoApi.infoApplicationMecPlatformGet(1, (error, data, response) => {
          if (error != null) {
            console.log(error);
          } else {
            setAppInfo(response.body);
            console.log(response.body);
          }
        });
      }, 1000);
      return () => clearInterval(interval);
    }
  }, [appInfo, registered]);

  // Peform polling activity logs only if log changes or if it registered
  useEffect(() => {
    if (registered) {
      const interval = setInterval(() => {
        appInfoApi.infoLogsGet(20, (error, data, response) => {
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

  // Stop polling app info by switching registeration into false + clean app info + one request to app info + one request to ams
  const deRegisterapp = async () => {
    setAppInfo({});
    console.log("hello world");
    setRegisteration(false);
    appInfoApi.infoApplicationMecPlatformDeleteDelete(1, (err, data, resp) => {
      if (err != null) {
        console.log(err);
      }
      appInfoApi.infoLogsGet(20, (error, data, response) => {
        if (error != null) {
          console.log(error);
        } else {
          console.log("reset");
          console.log(response.body);
          setAppLog(response.body);
        }
      });
      appInfoApi.infoApplicationMecPlatformGet(1, (error, data, response) => {
        if (error != null) {
          console.log(error);
        } else {
          setAppInfo(response.body);
          console.log(appInfo);
        }
      });
    });
  };

  // Stop polling app info by switching registeration into false + clean app info + one request to app info + one request to ams
  const removeAmsDevice = async (device) => {
    appInfoApi.serviceAmsDeleteDeviceDelete(device, (err, data, resp) => {
      if (err != null) {
        console.log(err);
      }
      console.log("byeeeee");
    });
  };

  return (
    <div className="ui-background" style={{ height: "100%" }}>
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
                style={{
                  backgroundColor: "#379DD8",
                  color: "white",
                }}
                fullWidth
                variant="contained"
                sx={{ mt: 3, mb: 2 }}
                onClick={() => {
                  appInfoApi.registerAppMecPlatformPost(
                    1,
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
                style={{
                  backgroundColor: "#379DD8",
                  color: "white",
                }}
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

            <div style={{ marginTop: "5rem" }}>
              <Button
                outlined
                style={{
                  backgroundColor: "#379DD8",
                  color: "white",
                }}
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
              <DialogContent>
                Create a new application mobility service resource
              </DialogContent>
              <TextField
                style={{ margin: "10px" }}
                label="device"
                value={textValue}
                onChange={(event) => {
                  console.log(event.target.value);
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
                    console.log(textValue);
                    appInfoApi.serviceAmsUpdateDevicePut(
                      textValue,
                      (err, data, resp) => {
                        if (err != null) {
                          console.log(err);
                        }
                        console.log(resp);
                      }
                    );
                  }}
                >
                  Confirm
                </DialogButton>
              </DialogActions>
            </Dialog>
            <Button
              outlined
              style={{
                backgroundColor: "#379DD8",
                color: "white",
              }}
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
                style={{ margin: "10px" }}
                label="device"
                value={textValue}
                onChange={(event) => {
                  console.log(event.target.value);
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
            <AmsPane data={amsLog} style={{ marginTop: "5px" }}></AmsPane>
          </GridCell>
        </Grid>
      </Container>
    </div>
  );
}

const styles = {
  headlineGrid: {
    marginBottom: 10,
  },
  title: {
    paddingTop: 10,
    paddingBottom: 15,
  },
};
