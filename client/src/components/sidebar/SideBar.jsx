import { useState } from 'react';
import { NavLink } from 'react-router-dom';
import inventory from '../../assets/inventory.svg';
import dashboard from '../../assets/dashboard.svg';
import shipment from '../../assets/shipment.svg';
import order from '../../assets/order.svg';
import analytics from '../../assets/analytics.svg';
import support from '../../assets/support.svg';
import setting from '../../assets/setting.svg';
import arrowLeft from '../../assets/arrowLeft.svg';

export default function SideBar() {
  const [isCollapsed, setIsCollapsed] = useState(false);

  const sections = [
    { name: 'Dashboard', icon: dashboard, path: '/dashboard' },
    { name: 'Inventory', icon: inventory, path: '/inventory' },
    { name: 'Shipments', icon: shipment, path: '/shipment' },
    { name: 'Orders', icon: order, path: '/order' },
    { name: 'Analytics', icon: analytics, path: '/analytics' },
    { name: 'Support', icon: support, path: '/support' },
    { name: 'Settings', icon: setting, path: '/setting' },
  ];

  return (
    <div
      className={`${
        isCollapsed ? 'w-20' : 'w-52'
      } bg-white dark:bg-gray-800 h-screen shadow-lg transition-all duration-500 overflow-hidden`}
    >
      {/* Toggle Button */}
      <div className="flex justify-end p-2">
        <button
          className="text-gray-600 dark:text-gray-300"
          onClick={() => setIsCollapsed(!isCollapsed)}
        >
          <img
            src={arrowLeft}
            alt="Toggle Arrow"
            className={`h-6 w-6 transform transition-transform duration-300 dark:invert ${
              isCollapsed ? 'rotate-180' : ''
            }`}
          />
        </button>
      </div>

      {/* Navigation Links */}
      <nav className="flex flex-col space-y-2 p-2">
        {sections.map((section) => (
          <NavLink
            key={section.name}
            to={section.path}
            className={({ isActive }) =>
              `group flex items-center gap-3 p-2 rounded hover:bg-gray-100 dark:hover:bg-gray-900 hover: transition-all ${
                isActive ? 'dark:bg-gray-900 bg-gray-200' : ''
              }`
            }
          >
            {/* Icon */}
            <div className="flex items-center justify-center h-10 w-10 flex-shrink-0">
              <img
                src={section.icon}
                alt={`${section.name} icon`}
                className="h-6 w-6 dark:invert group-hover:scale-125 transition-transform duration-300"
              />
            </div>

            {/* Text */}

            <span
              className={`text-sm font-medium transition-all duration-300 ease-in-out group-hover:scale-125 ${
                isCollapsed ? 'opacity-0 scale-90' : 'opacity-100 scale-100'
              }`}
            >
              {section.name}
            </span>
          </NavLink>
        ))}
      </nav>
    </div>
  );
}
