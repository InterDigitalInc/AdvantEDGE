import React from "react";
import { Typography } from "@rmwc/typography";
import "@material/typography/dist/mdc.typography.css";
import { Elevation } from "@rmwc/elevation";
// import { Grid, GridCell } from "@rmwc/grid";
import Title from "@/js/components/Title";
import Grid from "@mui/material/Grid";
import Paper from "@mui/material/Paper";
//import Typography from "@mui/material/Typography";
import Link from "@mui/material/Link";

// import css
import "@/css/pane.module.css";
import { element } from "prop-types";

export default function Activitypane({ height, data }) {
  function discoveredService(services) {
    if (services) {
      const resp = services.map((element) => {
        return (
          // <Typography style={{ display: "block" }}>
          //   {element.serName}
          //   &nbsp;&nbsp;&nbsp;
          //   {element.serInstanceId}
          //   &nbsp;&nbsp;&nbsp;
          //   {element.version}
          //   &nbsp;&nbsp;&nbsp;
          //   <Link href={element.link}> {" -- Link "}</Link>

          <Grid container>
            <Grid Paper xs={3}>
              {element.serName}
            </Grid>
            <Grid Paper xs={6}>
              {` Id: `}
              {element.serInstanceId}
            </Grid>
            <Grid Paper xs={1}>
              {element.version}
            </Grid>
            <Grid Paper xs={2}>
              <Link href={element.link}> {" -- Link "}</Link>
            </Grid>
          </Grid>

          /*       {/* <GridCell span="3">{element.serName}</GridCell>

            <GridCell span="6" style={{ marginBottom: "-5px" }}>
              {` Id: `}
              {element.serInstanceId}
            </GridCell>
            <GridCell span="1"> {element.version}</GridCell>
            <GridCell span="2">
            
            </GridCell> */
          /* </Typography> */
        );
      });
      if (resp) return resp;
    }
  }

  // TODO: Add logic for termination
  const computeData = (data) => {
    return (
      <div
        style={{
          height: "100vh",
          display: "block",
          wordWrap: "break-word",
          overflowY: "auto",
        }}
      >
        <Title>Application Instance</Title>
        <Typography theme="primary" use="subtitle2">
          MEC Platform
        </Typography>
        <div style={{ marginTop: "0.0.1rem" }}>
          <Typography use="caption">
            Name: {data?.name ? data.name : "N/A"}
          </Typography>
        </div>

        <Typography use="caption">
          Url: {data?.url ? data.url : "N/A"}
        </Typography>
        <div style={{ marginTop: "0.1rem" }}>
          <Typography theme="primary" use="subtitle2">
            Instance Info
          </Typography>
        </div>
        <div style={{ marginTop: "0.1rem" }}>
          <Typography use="caption" style={{ display: "block" }}>
            Config: {data?.config ? data.config : "N/A"}
          </Typography>
          <Typography use="caption" style={{ display: "block" }}>
            {`Ip & Port:`} {data?.ip ? data.ip : "N/A"}
          </Typography>
          <Typography use="caption" style={{ display: "block" }}>
            Id: {data?.id ? data.id : "N/A"}
          </Typography>
        </div>
        <div>
          <Typography use="caption">
            MEC011 Ready {`: `} {data?.mecReady ? "True" : "False"}
          </Typography>
          <Typography use="caption" style={{ display: "block" }}>
            MEC011 Terminated {`: `}
            {data?.mecTerminated ? "True" : "False"}
          </Typography>
          <Typography use="caption">
            MEC021 AMS resource {`: `} {data?.amsResource ? "True" : "False"}
          </Typography>
          <Typography use="caption" style={{ display: "block" }}>
            Subscriptions: {``}
          </Typography>
        </div>
        <div style={{ marginLeft: "2rem", marginTop: "0.1rem" }}>
          <Typography use="caption" style={{ display: "block" }}>
            {`Termination Id: `}
            {data?.subscriptions?.AppTerminationSubscription?.subId
              ? data?.subscriptions.AppTerminationSubscription?.subId
              : " N/A"}
          </Typography>
          <Typography use="caption" style={{ display: "block" }}>
            {`Service Availability Id: `}
            {data?.subscriptions?.SerAvailabilitySubscription?.subId
              ? data?.subscriptions.SerAvailabilitySubscription?.subId
              : "N/A"}
          </Typography>
          <Typography use="caption" style={{ display: "block" }}>
            {`Application Mobility Id: `}
            {data?.subscriptions?.AmsLinkListSubscription?.subId
              ? data?.subscriptions.AmsLinkListSubscription?.subId
              : "N/A"}
          </Typography>
        </div>
        <div style={{ marginTop: "0.1rem" }}>
          <Typography
            style={{ display: "block" }}
            theme="primary"
            use="subtitle2"
          >
            Offered Service
          </Typography>

          <Typography use="caption">
            Service Name:{" "}
            {data?.offeredService?.serName
              ? data.offeredService.serName
              : "N/A"}
          </Typography>

          <Typography use="caption" style={{ display: "block" }}>
            Id: {data?.offeredService?.id ? data?.offeredService?.id : "N/A"}
          </Typography>
          <Typography use="caption" style={{ display: "block" }}>
            State:{" "}
            {data?.offeredService?.state ? data.offeredService.state : "N/A"}
          </Typography>
          <Typography use="caption" style={{ display: "block" }}>
            Scope of Locality {`: `}{" "}
            {data?.offeredService?.scopeOfLocality
              ? data?.offeredService?.scopeOfLocality
              : "N/A"}
          </Typography>
          <Typography use="caption" style={{ display: "block" }}>
            Consumed Local Only {`: `}{" "}
            {data?.offeredService?.consumedLocalOnly ? "True" : "N/A"}
          </Typography>
        </div>
        <div style={{ marginTop: "0.1rem" }}>
          <Typography theme="primary" use="subtitle2">
            Discovered Service:
          </Typography>
          <Typography use="caption">
            {discoveredService(data?.discoveredServices)}
          </Typography>
        </div>
      </div>
    );
  };

  return (
    <div style={{ backgroundColor: "ffffff" }}>
      <Elevation
        z={2}
        className="component-style "
        style={{ padding: 10, marginBottom: 10 }}
      >
        {computeData(data)}
      </Elevation>
    </div>
  );
}
