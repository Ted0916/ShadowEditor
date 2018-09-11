﻿import UI from '../../ui/UI';
import BasicComponent from '../../component/BasicComponent';
import TransformComponent from '../../component/TransformComponent';
import CameraComponent from '../../component/CameraComponent';
import LightComponent from '../../component/LightComponent';
import ShadowComponent from '../../component/ShadowComponent';
import GeometryComponent from '../../component/GeometryComponent';
import MaterialComponent from '../../component/MaterialComponent';
import SceneComponent from '../../component/SceneComponent';
import AudioListenerComponent from '../../component/audio/AudioListenerComponent';

/**
 * 属性面板
 * @author mrdoob / http://mrdoob.com/
 * @author tengge / https://github.com/tengge1
 */
function PropertyPanel(options) {
    UI.Control.call(this, options);
    this.app = options.app;
};

PropertyPanel.prototype = Object.create(UI.Control.prototype);
PropertyPanel.prototype.constructor = PropertyPanel;

PropertyPanel.prototype.render = function () {
    var data = {
        xtype: 'div',
        parent: this.parent,
        children: [
            new BasicComponent({ app: this.app }),
            new TransformComponent({ app: this.app }),
            new SceneComponent({ app: this.app }),
            new CameraComponent({ app: this.app }),
            new LightComponent({ app: this.app }),
            new ShadowComponent({ app: this.app }),
            new GeometryComponent({ app: this.app }),
            new MaterialComponent({ app: this.app }),
            new AudioListenerComponent({ app: this.app })
        ]
    };

    var control = UI.create(data);
    control.render();
};

export default PropertyPanel;