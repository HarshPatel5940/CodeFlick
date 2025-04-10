//  This file is not required in the project but it is required to be present for the TailwindCSS IntelliSense to work.
import type { Config } from "tailwindcss";

export default <Partial<Config>>{
  theme: {
    extend: {
      colors: {
        mybg: '#061326',
        myborder: '#2c3c54',
        myLightBg: '#F4F5F9',
        myLightBorder: '#282626',
      }
    }
  }
};
