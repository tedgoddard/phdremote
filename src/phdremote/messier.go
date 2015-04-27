package phdremote

//http://messier.seds.org/dataRA.html
//awk '{dm = int(10000 * $9 / 60); printf("\"%s\": {%g,%s.%04.f},\n", $1, (($6 + $7 / 60.0 ) * 15.0), $8, dm)}'
var Messier = map[string][]float32 {
    "M110": {10.1,+41.6833},
    "M31": {10.675,+41.2666},
    "M32": {10.675,+40.8666},
    "M103": {23.3,+60.7000},
    "M33": {23.475,+30.6500},
    "M74": {24.175,+15.7833},
    "M76": {25.6,+51.5666},
    "M34": {40.5,+42.7833},
    "M77": {40.675,-00.0166},
    "M45": {56.75,+24.1166},
    "M79": {81.05,-24.5166},
    "M38": {82.1,+35.8333},
    "M1": {83.625,+22.0166},
    "M42": {83.85,-05.4500},
    "M43": {83.9,-05.2666},
    "M36": {84.025,+34.1333},
    "M78": {86.675,+00.0500},
    "M37": {88.1,+32.5500},
    "M35": {92.225,+24.3333},
    "M41": {101.5,-20.7333},
    "M50": {105.8,-08.3333},
    "M47": {114.15,-14.5000},
    "M46": {115.45,-14.8166},
    "M93": {116.15,-23.8666},
    "M48": {123.45,-05.8000},
    "M44": {130.025,+19.9833},
    "M67": {132.6,+11.8166},
    "M81": {148.9,+69.0666},
    "M82": {148.95,+69.6833},
    "M95": {161,+11.7000},
    "M96": {161.7,+11.8166},
    "M105": {161.95,+12.5833},
    "M108": {167.875,+55.6666},
    "M97": {168.7,+55.0166},
    "M65": {169.725,+13.0833},
    "M66": {170.05,+12.9833},
    "M109B": {178.45,+52.3333},
    "M109": {179.4,+53.3833},
    "M98": {183.45,+14.9000},
    "M99": {184.7,+14.4166},
    "M106": {184.75,+47.3000},
    "M61": {185.475,+04.4666},
    "M40": {185.6,+58.0833},
    "M100": {185.725,+15.8166},
    "M84": {186.275,+12.8833},
    "M85": {186.35,+18.1833},
    "M86": {186.55,+12.9500},
    "M49": {187.45,+08.0000},
    "M87": {187.7,+12.4000},
    "M88": {188,+14.4166},
    "M91": {188.85,+14.5000},
    "M89": {188.925,+12.5500},
    "M90": {189.2,+13.1666},
    "M58": {189.425,+11.8166},
    "M68": {189.875,-26.7500},
    "M104": {190,-11.6166},
    "M59": {190.5,+11.6500},
    "M60": {190.925,+11.5500},
    "M94": {192.725,+41.1166},
    "M64": {194.175,+21.6833},
    "M53": {198.225,+18.1666},
    "M63": {198.95,+42.0333},
    "M51": {202.475,+47.2000},
    "M51B": {202.5,+47.2666},
    "M83": {204.25,-29.8666},
    "M3": {205.55,+28.3833},
    "M101": {210.8,+54.3500},
    "M102": {226.625,+55.7666},
    "M5": {229.65,+02.0833},
    "M80": {244.25,-22.9833},
    "M4": {245.9,-26.5333},
    "M107": {248.125,-13.0500},
    "M13": {250.425,+36.4666},
    "M12": {251.8,-01.9500},
    "M10": {254.275,-04.1000},
    "M62": {255.3,-30.1166},
    "M19": {255.65,-26.2666},
    "M92": {259.275,+43.1333},
    "M9": {259.8,-18.5166},
    "M14": {264.4,-03.2500},
    "M6": {265.025,-32.2166},
    "M7": {268.475,-34.8166},
    "M23": {269.2,-19.0166},
    "M20": {270.65,-23.0333},
    "M8": {270.95,-24.3833},
    "M21": {271.15,-22.5000},
    "M24": {274.225,-18.4833},
    "M16": {274.7,-13.7833},
    "M18": {274.975,-17.1333},
    "M17": {275.2,-16.1833},
    "M28": {276.125,-24.8666},
    "M69": {277.85,-32.3500},
    "M25": {277.9,-19.2500},
    "M22": {279.1,-23.9000},
    "M70": {280.8,-32.3000},
    "M26": {281.3,-09.4000},
    "M11": {282.775,-06.2666},
    "M57": {283.4,+33.0333},
    "M54": {283.775,-30.4833},
    "M56": {289.15,+30.1833},
    "M55": {295,-30.9666},
    "M71": {298.45,+18.7833},
    "M27": {299.9,+22.7166},
    "M75": {301.525,-21.9166},
    "M29": {305.975,+38.5333},
    "M72": {313.375,-12.5333},
    "M73": {314.725,-12.6333},
    "M15": {322.5,+12.1666},
    "M39": {323.05,+48.4333},
    "M2": {323.375,-00.8166},
    "M30": {325.1,-23.1833},
    "M52": {351.05,+61.5833},
}
