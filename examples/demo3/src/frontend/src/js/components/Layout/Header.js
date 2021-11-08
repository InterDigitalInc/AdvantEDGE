import React from "react";
import {
  Toolbar,
  ToolbarRow,
  ToolbarSection,
  ToolbarTitle,
} from "@rmwc/toolbar";
import { Elevation } from "@rmwc/elevation";
import "@material/toolbar/dist/mdc.toolbar.css";
import "@material/elevation/dist/mdc.elevation.css";

export default function header_feat() {
  const logo = require("@/img/ID-Icon-01-idcc.svg");
  const advantEdge = require("@/img/AdvantEDGE-logo-NoTagline_White_RGB.png");
  return (
    <div style={{ height: "48px" }}>
      <Toolbar fixed style={{ zIndex: 9000, backgroundColor: "#379DD8" }}>
        <Elevation z={2}>
          <ToolbarRow>
            <ToolbarSection alignStart style={{ display: "contents" }}>
              <div style={styles.flex}>
                <img
                  src={logo}
                  alt=""
                  style={{
                    height: "20px",
                    marginLeft: "10px",
                  }}
                />
                <img height={50} src={advantEdge} alt="" />
                <ToolbarTitle>Demo 3 MEC Edge Application</ToolbarTitle>
              </div>
            </ToolbarSection>
          </ToolbarRow>
        </Elevation>
      </Toolbar>
    </div>
  );
}

const styles = {
  flex: {
    display: "flex",
    justifyContent: "center",
    alignItems: "center",
  },
};
